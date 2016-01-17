# drone-slack

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-slack/status.svg)](http://beta.drone.io/drone-plugins/drone-slack)
[![](https://badge.imagelayers.io/plugins/drone-slack:latest.svg)](https://imagelayers.io/?images=plugins/drone-slack:latest 'Get your own badge on imagelayers.io')

Drone plugin for sending Slack notifications

## Overview

This plugin is responsible for sending build notifications to your Slack channel:

```sh
./drone-slack <<EOF
{
    "repo" : {
        "owner": "foo",
        "name": "bar",
        "full_name": "foo/bar"
    },
    "system": {
        "link_url": "http://drone.mycompany.com"
    },
    "build" : {
        "number": 22,
        "status": "success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "commit": "9f2849d5",
        "branch": "master",
        "message": "Update the Readme",
        "author": "johnsmith",
        "author_email": "john.smith@gmail.com"
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
        "owner": "foo",
        "name": "bar",
        "full_name": "foo/bar"
    },
    "system": {
        "link_url": "http://drone.mycompany.com"
    },
    "build" : {
        "number": 22,
        "status": "success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "commit": "9f2849d5",
        "branch": "master",
        "message": "Update the Readme",
        "author": "johnsmith",
        "author_email": "john.smith@gmail.com"
    },
    "vargs": {
        "webhook_url": "https://hooks.slack.com/services/...",
        "username": "drone",
        "channel": "#dev"
    }
}
EOF
```
