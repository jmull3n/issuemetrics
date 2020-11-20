package issuemetrics

import (
	"github.com/jmull3n/issuemetrics/pkg/metrics"
	"sort"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// Munger is the place where the munging happens
type Munger struct {
	// internal processing bits
	issues         []*Issue
	repos          []string
	topDaySummary  map[time.Time]int
	repoDaySummary map[time.Time]map[string]int
	errors         []error

	// locks so we can make the munger threadsafe
	topdaylock  sync.Mutex
	repodaylock sync.Mutex
	issueslock  sync.Mutex
	reposlock   sync.Mutex
	errorslock  sync.Mutex
}

// NewMunger returns a new munger reference
func NewMunger() *Munger {
	return &Munger{
		issues:         []*Issue{},
		topDaySummary:  make(map[time.Time]int),
		repoDaySummary: make(map[time.Time]map[string]int),
	}
}

// SetRepo is a threadsafe setter for a repo
func (m *Munger) SetRepo(repo string) {
	m.reposlock.Lock()
	defer m.reposlock.Unlock()

	m.repos = append(m.repos, repo)
}

// GetRepos is a threadsafe getter for a repo list
func (m *Munger) GetRepos() []string {
	m.reposlock.Lock()
	defer m.reposlock.Unlock()
	sort.Slice(m.repos, func(i, j int) bool {
		return m.repos[i] < m.repos[j]
	})
	return m.repos
}

// AddError adds an error to the munger
func (m *Munger) AddError(err error) {
	m.errorslock.Lock()
	defer m.errorslock.Unlock()
	m.errors = append(m.errors, err)
}

// GetErrors gets the errors
func (m *Munger) GetErrors() []error {
	m.errorslock.Lock()
	defer m.errorslock.Unlock()
	return m.errors
}

// GetIssues gets the munger's issues
func (m *Munger) GetIssues() []*Issue {
	m.issueslock.Lock()
	defer m.issueslock.Unlock()

	sort.Slice(m.issues, func(i, j int) bool {
		return m.issues[i].CreatedAt.Before(m.issues[j].CreatedAt)
	})

	return m.issues
}

// ProcessIssues performs a threadsafe add of issues to the munging object
func (m *Munger) ProcessIssues(issues []*Issue) {
	ts := time.Now()
	m.issueslock.Lock()
	defer m.issueslock.Unlock()

	// go through the incoming issues and add the bits to the munger state
	for _, issue := range issues {
		// update the repo issue/day grouping
		m.incrRepoDayIssue(issue.CreatedAt, issue.Repo)

		// update the top issue day grouping
		m.updateTopDay(issue.CreatedAt)
	}
	// sample the duration so we can see if this is really slow
	metrics.RequestDurationMilliseconds.WithLabelValues("issuemetrics", "munger:ProcessIssues", "200").Observe(float64(time.Since(ts).Milliseconds()))
	m.issues = append(m.issues, issues...)
}

// incrRepoDayIssue performs a threadsafe update of the repoday issuecount
func (m *Munger) incrRepoDayIssue(createdAt time.Time, repo string) {
	m.repodaylock.Lock()
	defer m.repodaylock.Unlock()
	day := createdAt.Truncate(time.Hour * 24)
	// check to see if the day is in the collection and create a new map for the repo
	_, ok := m.repoDaySummary[day]
	if !ok {
		m.repoDaySummary[day] = make(map[string]int)
	}

	// set the repo summary on the day
	m.repoDaySummary[day][repo]++
}

// updateTopDay performs a threadsafe update of the topday summary
func (m *Munger) updateTopDay(createdAt time.Time) {
	m.topdaylock.Lock()
	defer m.topdaylock.Unlock()
	day := createdAt.Truncate(time.Hour * 24)

	// check to see if the day is in the collection and create the row or add to it
	m.topDaySummary[day]++
	log.Debug(m.topDaySummary)

}

// GetTopDay gets the top day from the topDaySummary collection
func (m *Munger) GetTopDay() (*time.Time, map[string]int) {
	ts := time.Now()

	m.topdaylock.Lock()
	defer m.topdaylock.Unlock()

	keys := make([]time.Time, 0, len(m.topDaySummary))
	for key := range m.topDaySummary {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		// sort by value
		if m.topDaySummary[keys[i]] != m.topDaySummary[keys[j]] {
			return m.topDaySummary[keys[i]] > m.topDaySummary[keys[j]]
		}

		// if 2 days have the same value, then sort by date
		return keys[i].After(keys[j])

	})
	// since the summay keys are sorted in desc order, just returning the first will get us the right day
	result := keys[0]

	// sample the duration so we can see if this is really slow
	metrics.RequestDurationMilliseconds.WithLabelValues("issuemetrics", "munger:GetTopDay", "200").Observe(float64(time.Since(ts).Milliseconds()))
	return &result, m.repoDaySummary[keys[0]]
}
