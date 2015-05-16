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
        "name": "bar",
        "self_url": "http://my.drone.io/foo/bar"
    },
    "commit" : {
        "state": "success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "sha": "9f2849d5",
        "branch": "master",
        "pull_request": "800",
        "author": "john.smith@gmail.com",
        "message": "Update the Readme"
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

Build the Docker container. Note that we need to use the `-netgo` tag so that
the binary is built without a CGO dependency:

```sh
CGO_ENABLED=0 go build -a -tags netgo
docker build --rm=true -t plugins/drone-slack .
```

Send a Slack notification:

```sh
docker run -i plugins/drone-slack <<EOF
{
    "repo" : {
        "host": "github.com",
        "owner": "foo",
        "name": "bar",
        "self": "http://my.drone.io/foo/bar"
    },
    "commit" : {
        "state": "success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "sha": "9f2849d5",
        "branch": "master",
        "pull_request": "800",
        "author": "john.smith@gmail.com",
        "message": "Update the Readme"
    },
    "vargs": {
        "webhook_url": "https://hooks.slack.com/services/...",
        "username": "drone", 
        "channel": "#dev"
    }
}
EOF
```
