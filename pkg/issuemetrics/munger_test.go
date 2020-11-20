package issuemetrics

import (
	"fmt"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMarshalIssue(t *testing.T) {
	rawissues := `
		[{
		  "id": 38,
		  "state": "open",
		  "title": "Found a bug",
		  "repository": "owner1/repository1",
		  "created_at": "2011-04-22T13:33:48Z"
		}]`
	var issues []Issue
	err := json.Unmarshal([]byte(rawissues), &issues)

	assert.Nil(t, err)

	assert.Equal(t, 38, issues[0].ID)
	assert.Equal(t, "open", issues[0].State)
	assert.Equal(t, "owner1/repository1", issues[0].Repo)
	assert.Equal(t, "2011-04-22T13:33:48Z", issues[0].CreatedAt.Format(time.RFC3339))

}

func TestProcessIssues(t *testing.T) {
	assert := assert.New(t)
	rawissues := `
	[
		{
		  "id": 38,
		  "state": "open",
		  "title": "Found a bug",
		  "repository": "owner1/repository1",
		  "created_at": "2011-04-22T13:33:48Z"
		},
		{
		  "id": 23,
		  "state": "open",
		  "title": "Found a bug 2",
		  "repository": "owner1/repository1",
		  "created_at": "2011-04-22T18:24:32Z"
		},
		{
		  "id": 24,
		  "state": "closed",
		  "title": "Feature request",
		  "repository": "owner2/repository2",
		  "created_at": "2011-05-08T09:15:20Z"
		}
	  ]`
	var issues []*Issue
	json.Unmarshal([]byte(rawissues), &issues)

	munger := NewMunger()

	munger.ProcessIssues(issues)
	ok := false
    if munger.issues[0].ID == 38 {
		ok = true
	}

	fmt.Printf("%v", munger.topDaySummary)
	assert.True(ok)

	ok = false
	topDay, occurances := munger.GetTopDay()
	fmt.Printf("%v: %v", topDay, occurances)

	if topDay.Format(time.RFC3339) == "2011-04-22T00:00:00Z"{
		ok = true
	}
	assert.True(ok)
}

