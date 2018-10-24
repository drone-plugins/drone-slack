# drone-slack

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-slack/status.svg)](http://beta.drone.io/drone-plugins/drone-slack)
[![Join the discussion at https://discourse.drone.io](https://img.shields.io/badge/discourse-forum-orange.svg)](https://discourse.drone.io)
[![Drone questions at https://stackoverflow.com](https://img.shields.io/badge/drone-stackoverflow-orange.svg)](https://stackoverflow.com/questions/tagged/drone.io)
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-slack?status.svg)](http://godoc.org/github.com/drone-plugins/drone-slack)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-slack)](https://goreportcard.com/report/github.com/drone-plugins/drone-slack)
[![](https://images.microbadger.com/badges/image/plugins/slack.svg)](https://microbadger.com/images/plugins/slack "Get your own image badge on microbadger.com")

Drone plugin for sending Slack notifications. For the usage information and a listing of the available options please take a look at [the docs](http://plugins.drone.io/drone-plugins/drone-slack/).

## Build

Build the binary with the following commands:

```
go build
```

## Docker

Build the Docker image with the following commands:

```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -o release/linux/amd64/drone-slack
docker build --rm -t plugins/slack .
```

## Usage

Execute from the working directory:

```
docker run --rm \
  -e SLACK_WEBHOOK=https://hooks.slack.com/services/... \
  -e PLUGIN_CHANNEL=foo \
  -e PLUGIN_USERNAME=drone \
  -e DRONE_REPO_OWNER=octocat \
  -e DRONE_REPO_NAME=hello-world \
  -e DRONE_COMMIT_SHA=7fd1a60b01f91b314f59955a4e4d4e80d8edf11d \
  -e DRONE_COMMIT_BRANCH=master \
  -e DRONE_COMMIT_AUTHOR=octocat \
  -e DRONE_BUILD_NUMBER=1 \
  -e DRONE_BUILD_STATUS=success \
  -e DRONE_BUILD_LINK=http://github.com/octocat/hello-world \
  -e DRONE_TAG=1.0.0 \
  plugins/slack
```

## S3 Name Translation

```
PLUGIN_S3_BUCKET
PLUGIN_S3_USERS_KEY
```

These environment variables/options can be set to an S3 bucket and file key respectively. This will
cause the plugin use the json map in that file to translate `build.author` and `recipient` fields
into `computed.authorSlack` and `computed.recipientSlack` leaving the original fields as-is.

These new fields can be used in templates normally. The recipient option (to send a slack message
to a specific user vs a channel) with automatically translate the name when the new S3 options are
set.

The format of the JSON file:
{
   "gitUserName": "slackUserName",
}

**Note** When testing this locally make sure to set the aws authentication environment variables on
the docker container.