# drone-slack

Drone plugin for sending Slack notifications

## Build

Build the binary with the following commands:

```
export GO15VENDOREXPERIMENT=1
go build
go test
```

## Docker

Build the docker image with the following commands:

```
export GO15VENDOREXPERIMENT=1
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo
```

Please note incorrectly building the image for the correct x64 linux and with GCO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-slack' not found or does not exist..
```

## Usage

Post the build status to a channel:

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
    plugins/slack
```
