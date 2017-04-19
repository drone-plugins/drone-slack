# drone-slack

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-slack/status.svg)](http://beta.drone.io/drone-plugins/drone-slack)
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-slack?status.svg)](http://godoc.org/github.com/drone-plugins/drone-slack)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-slack)](https://goreportcard.com/report/github.com/drone-plugins/drone-slack)
[![Join the chat at https://gitter.im/drone/drone](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/drone/drone)

Drone plugin for sending Slack notifications. For the usage information and a listing of the available options please reference [the documentation](http://plugins.drone.io/drone-plugins/slack).

## Build

Build the binary with the following commands:

```
go build
go test
```

## Docker

Build the docker image with the following commands:

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo
docker build -t plugins/slack .
```

Please note incorrectly building the image for the correct x64 linux and with
GCO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-slack' not found or does not exist..
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
