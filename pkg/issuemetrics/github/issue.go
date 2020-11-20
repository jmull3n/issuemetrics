package github

import (
	"time"
	"encoding/json"
	"fmt"
	"github.com/jmull3n/issuemetrics/pkg/issuemetrics"
	"github.com/jmull3n/issuemetrics/pkg/metrics"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

// IssueManager is the manager for the github issues
type IssueManager struct {
}

type githubIssue struct {
	issuemetrics.Issue
	RepoURL string `json:"repository_url"`
}

// enforce the contract
var _ issuemetrics.Manager = (*IssueManager)(nil)

func getIssues(repository string) ([]*issuemetrics.Issue, error) {
	urlString := fmt.Sprintf("https://api.github.com/repos/%s/issues", repository)
	response, err := http.Get(urlString)
	log.Info(response)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer response.Body.Close()
	var issues []*issuemetrics.Issue

	if err := json.NewDecoder(response.Body).Decode(&issues); err != nil {
		log.Error("error parsing issue payload", "err", err)
		return nil, err
	}

	// do the transformation for the respository, since github sends down a repository_url
	for _, issue := range issues {
		issue.Repo = repository
	}

	return issues, nil
}

func processIssuesForRepos(munger *issuemetrics.Munger, input issuemetrics.Input) {
	var wg sync.WaitGroup

	// make semaphore buffer so we can control the concurrency
	buffer := make(chan bool, 10) // TODO: make concurrency configurable

	for _, repo := range input.Repos {
		wg.Add(1)
		go func(repository string) {
			ts := time.Now()
			defer wg.Done() // wait for this to finish since it'll be done async per repo
			buffer <- true
			munger.SetRepo(repository)

			issues, err := getIssues(repository)
			statuscode := "200"
			if err != nil {
				statuscode = "500"
				log.Error(err, "repo:", repository)
				munger.AddError(err) // TODO: make this a little better, but lets see what hapens first
			}
			munger.ProcessIssues(issues)
			metrics.RequestDurationMilliseconds.WithLabelValues("issuemetrics", "github:processIssuesForRepos", statuscode).Observe(float64(time.Since(ts).Milliseconds()))
			<-buffer
		}(repo)
	}
	wg.Wait()
}

func makeOutput(munger *issuemetrics.Munger) *issuemetrics.Output {

	day, actualTopDayOccurances := munger.GetTopDay()

	occurances := make(map[string]int)

	for _, repo := range munger.GetRepos() {
		val, ok := actualTopDayOccurances[repo]
		if !ok {
			occurances[repo] = 0
		} else {
			occurances[repo] = val
		}
	}

	return &issuemetrics.Output{
		Issues: munger.GetIssues(),
		Day: issuemetrics.TopDay{
			Occurrances: occurances,
			Day:         day.Format("2006-01-02"),
		},
		Errors: munger.GetErrors(),
	}
}

// IssueMetricsHandler is the handler that does the github issue stuff
func (g *IssueManager) IssueMetricsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input issuemetrics.Input

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Error("error parsing input", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// get an instance of a munger for this request
	munger := issuemetrics.NewMunger()

	processIssuesForRepos(munger, input)

	output := makeOutput(munger)
	response, err := json.Marshal(output)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	size, err := w.Write(response)
	statuscode := "200"
	if err != nil {
		statuscode = "500"
		log.Error("error parsing input", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	// track output size
	metrics.ResponseBytes.WithLabelValues("issumetrics", "github", statuscode).Observe(float64(size))

}
