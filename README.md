

## Overview
The Issue metrics service accepts an input of:

```
{
    "repos": [
        "magiclabs/example-nextjs-faunadb-todomvc",
		"kedacore/keda"
    ]
}
```
and produces an output of the issues, sorted in descending order and the top day/repo(s) which issues were created.
```{
  "issues": [
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
  ],
  "top_day": {
    "day": "2011-04-22",
    "occurrences": {
      "owner1/repository1": 2,
      "owner2/repository2": 0
    }
  }
}
```

---
# Usage
Requirements:
- docker
- jq


```
make test

make run

curl -X POST -H "Content-Type: application/json" \
    -d '{ "repos": [ "magiclabs/example-nextjs-faunadb-todomvc", "kedacore/keda"]}' \
    localhost:8000/github | jq .
```

Navigate to `localhost:8083/metrics` to view the prometheus metrics that have been instrumented to get an intial understanding of how the service behaves during usage.

---
# Architecture
This is setup using a request/response pattern right now, but given that usage patterns may vary, a different strategy may be needed.

Right now, it's limited to the default 30 issues/repo of the github api, but under a realistic use-case, this may need to run on a schedule and catalog the issues in a dedicated storage layer.

