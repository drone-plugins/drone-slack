# drone-slack
Drone plugin for sending Slack notifications


## Overview

This plugin is responsible for sending build notifications to your Slack channel:

```sh
./drone-slack <<EOF
{
    "repo" : {
        "host": "github.com",
        "owner": "foo",
        "name": "bar"
    },
    "commit" : {
        "status": "Success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "sha": "9f2849d5",
        "branch": "master",
        "pull_request": "800",
        "author": "john.smith@gmail.com",
        "message": "Update the Readme"
    },
    "links" : {
        "repo_url": "https://foo.com/github.com/foo/bar",
        "commit_url": "https://foo.com/github.com/foo/bar/master/436b7a6e"
    },
    "vargs": {
        "webhook_url": "https://hooks.slack.com/services/...",
        "username": "drone", 
        "channel": "#dev"
    }
}
EOF
```

## Docker

Build the Docker container:

```sh
docker build -t drone-plugins/drone-slack .
```

Send a Slack notification:

```sh
docker run -i drone-plugins/drone-slack <<EOF
{
    "repo" : {
        "host": "github.com",
        "owner": "foo",
        "name": "bar"
    },
    "commit" : {
        "status": "Success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "sha": "9f2849d5",
        "branch": "master",
        "pull_request": "800",
        "author": "john.smith@gmail.com",
        "message": "Update the Readme"
    },
    "links" : {
        "repo_url": "https://foo.com/github.com/foo/bar",
        "commit_url": "https://foo.com/github.com/foo/bar/master/436b7a6e"
    },
    "vargs": {
        "webhook_url": "https://hooks.slack.com/services/...",
        "username": "drone", 
        "channel": "#dev"
    }
}
EOF
```
