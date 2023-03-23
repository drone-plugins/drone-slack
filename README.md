# drone-slack

[![Build Status](http://harness.drone.io/api/badges/drone-plugins/drone-slack/status.svg)](http://harness.drone.io/drone-plugins/drone-slack)
[![Slack](https://img.shields.io/badge/slack-drone-orange.svg?logo=slack)](https://join.slack.com/t/harnesscommunity/shared_invite/zt-y4hdqh7p-RVuEQyIl5Hcx4Ck8VCvzBw)
[![Join the discussion at https://community.harness.io](https://img.shields.io/badge/discourse-forum-orange.svg)](https://community.harness.io)
[![Drone questions at https://stackoverflow.com](https://img.shields.io/badge/drone-stackoverflow-orange.svg)](https://stackoverflow.com/questions/tagged/drone.io)
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-slack?status.svg)](http://godoc.org/github.com/drone-plugins/drone-slack)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-slack)](https://goreportcard.com/report/github.com/drone-plugins/drone-slack)

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
  -e DRONE_COMMIT_AUTHOR_EMAIL=octocat@github.com \
  -e DRONE_COMMIT_AUTHOR_AVATAR="https://avatars0.githubusercontent.com/u/583231?s=460&v=4" \
  -e DRONE_COMMIT_AUTHOR_NAME="The Octocat" \
  -e DRONE_BUILD_NUMBER=1 \
  -e DRONE_BUILD_STATUS=success \
  -e DRONE_BUILD_LINK=http://github.com/octocat/hello-world \
  -e DRONE_TAG=1.0.0 \
  plugins/slack
```

## Release Preparation

Run the changelog generator.

```BASH
docker run -it --rm -v "$(pwd)":/usr/local/src/your-app githubchangeloggenerator/github-changelog-generator -u drone-plugins -p drone-slack -t <secret github token>
```

You can generate a token by logging into your GitHub account and going to Settings -> Personal access tokens.

Next we tag the PR's with the fixes or enhancements labels. If the PR does not fufil the requirements, do not add a label.

**Before moving on make sure to update the version file `version/version.go && version/version_test.go`.**

Run the changelog generator again with the future version according to semver.

```BASH
docker run -it --rm -v "$(pwd)":/usr/local/src/your-app githubchangeloggenerator/github-changelog-generator -u drone-plugins -p drone-slack <secret token> --future-release v1.0.0
```

Create your pull request for the release. Get it merged then tag the release.
