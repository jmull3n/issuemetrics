package issuemetrics

import (
	"net/http"
	"sort"
	"time"
)

// Issue represents an issue
type Issue struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	Repo      string    `json:"repository"`
	CreatedAt time.Time `json:"created_at"`
}

// TopDay is the top day of the repos
type TopDay struct {
	Occurrances map[string]int `json:"occurances"`
	Day         string         `json:"day"`
}

// Input is the input datastructure of the api
type Input struct {
	// Repos in a {org}/{repo} string
	Repos []string `json:"repos"`
}

// Output is the datastructure of the issuemetrics output
type Output struct {
	Issues []*Issue `json:"issues"`
	Day    TopDay  `json:"top_day"`
	Errors []error `json:"errors"`
}

// SortIssuesAsc sorts the issues in ascending createdat
func (o *Output) SortIssuesAsc() {
	sort.Slice(o.Issues, func(i, j int) bool {
		return o.Issues[i].CreatedAt.Before(o.Issues[j].CreatedAt)
	})
}

// Manager is a interface for managing metric mungers
type Manager interface {
	// IssueMetricsHandler is a valid http handler
	IssueMetricsHandler(w http.ResponseWriter, r *http.Request)
}
