# Woodpecker CI Slack Plugin

Woodpecker plugin for sending Slack notifications cloned from: https://github.com/drone-plugins/drone-slack

For the usage information and a listing of the available options please take a look at [the docs](http://plugins.drone.io/drone-plugins/drone-slack/).

## Build

Build the binary with the following commands:

```
mise run build
```

## Docker

Build the Docker image with the following commands:

```
docker build --rm -t gowerstreet/slack -f docker/Dockerfile.linux.amd64 .
```

## Send Slack messages Usage

To Send Slack messages use the following

Execute from the working directory:

```
docker run --rm \
  -e SLACK_WEBHOOK=$SLACK_WEBHOOK \
  -e PLUGIN_CHANNEL=foo \
  -e PLUGIN_USERNAME=drone \
  -e CI_REPO_OWNER=octocat \
  -e CI_REPO_NAME=hello-world \
  -e CI_COMMIT_SHA=7fd1a60b01f91b314f59955a4e4d4e80d8edf11d \
  -e CI_COMMIT_BRANCH=master \
  -e CI_COMMIT_AUTHOR=octocat \
  -e CI_COMMIT_AUTHOR_EMAIL=octocat@github.com \
  -e CI_COMMIT_AUTHOR_AVATAR="https://avatars0.githubusercontent.com/u/583231?s=460&v=4" \
  -e CI_COMMIT_AUTHOR_NAME="The Octocat" \
  -e CI_COMMIT_TAG=1.0.0 \
  -e CI_PIPELINE_NUMBER=1 \
  -e CI_PIPELINE_STATUS=success \
  -e CI_PIPELINE_LINK=http://github.com/octocat/hello-world \
  gowerstreet/slack:local
```

Please note the following new environment variables:

- `SLACK_ACCESS_TOKEN`: The access token for Slack API authentication.
- `PLUGIN_CUSTOM_BLOCK`: Custom blocks in JSON format to include in the Slack message.

Make sure to replace `your_access_token` with your actual Slack access token and adjust

If you provide an access token, it will use the Slack API to send the message. Otherwise, it will use the webhook.

## Upload files to Slack

To Send Slack messages use the following

Execute from the working directory:
```
docker run --network host --rm \
  -e PLUGIN_ACCESS_TOKEN=your_access_token   \
  -e PLUGIN_CHANNEL=C07TL1KNV8Q \
  -e PLUGIN_USERNAME=jenkinstest003app \
  -e PLUGIN_FILE_PATH='/home/hns/test/b.txt' \
  -e PLUGIN_INITIAL_COMMENT='some start of text' \
  -e PLUGIN_TITLE='Build OK now' \
  plugins/slack
```

Please note the following new environment variables:

- `PLUGIN_ACCESS_TOKEN`: The access token for Slack API authentication.
- `PLUGIN_CUSTOM_BLOCK`: Custom blocks in JSON format to include in the Slack message.

Make sure to replace `your_access_token` with your actual Slack access token and adjust

If you provide an access token, it will use the Slack API to send the message.


### Get Slack Id of a user from a Email ID
```bash
docker run --network host --rm \
-e PLUGIN_ACCESS_TOKEN=your_access_token \
-e PLUGIN_SLACK_ID_OF=user01@somedomain.com \
plugins/slack
```
Output will be stored in the FOUND_SLACK_ID environment variable
Make sure to replace `your_access_token` with your actual Slack access token and adjust

### Get the Slack IDs of all committers from a git repo with two commit ids as commit Ids
```bash
docker run --network host --rm \
-e PLUGIN_ACCESS_TOKEN=your_access_token \
-e PLUGIN_COMMITTER_LIST_GIT_PATH=/harness \
-e PLUGIN_RECENT_COMMIT_ID=7fd1a60b01f91b314f59955a4e4d4e80d8edf11d \
-e PLUGIN_OLD_COMMIT_ID=7fd1a60b01f91b314f59955a4e4d4e80d8edf11d \
plugins/slack
```

### Get the Slack IDs of all committers from a git repo with HEAD and a commit id
```bash
docker run --network host --rm \
-e PLUGIN_ACCESS_TOKEN=your_access_token \
-e PLUGIN_COMMITTER_LIST_GIT_PATH=/harness \
-e PLUGIN_RECENT_COMMIT_ID=HEAD \
-e PLUGIN_OLD_COMMIT_ID=7fd1a60b01f91b314f59955a4e4d4e80d8edf11d \
plugins/slack
```

### Get the Slack IDs of all committers from a git repo with HEAD and a number of commits behind HEAD
```bash
docker run --network host --rm \
-e PLUGIN_ACCESS_TOKEN=your_access_token \
-e PLUGIN_COMMITTER_LIST_GIT_PATH=/harness \
-e PLUGIN_RECENT_COMMIT_ID=HEAD \
-e PLUGIN_OLD_COMMIT_ID=5 \
plugins/slack
```

Output will be stored in the COMMITTER_SLACK_ID_LIST environment variable as comma separated values.
Make sure to replace `your_access_token` with your actual Slack access token and adjust.

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
