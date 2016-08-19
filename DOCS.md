Use this plugin to send a message to your Slack channel when a build completes.
You will need to supply Drone with an [Incoming Webhook URL](https://my.slack.com/services/new/incoming-webhook).

## Config

The following parameters are used to configure the plugin:

* **webhook** - json payloads are sent here
* **channel** - messages sent to the above webhook are posted here
* **recipient** - alternatively you can send it to a specific user
* **username** - choose the username this integration will post as
* **template** - overwrite the default message template
* **image_url** - A valid URL to an image file that will be displayed inside a message attachment
* **icon_url** - A valid URL that displays a image to the left of the username
* **icon_emoji** - displays a emoji to the left of the username

The following secret values can be set to configure the plugin.

* **SLACK_WEBHOOK** - corresponds to **webhook**

It is highly recommended to put the **SLACK_WEBHOOK** into a secret so it is
not exposed to users. This can be done using the drone-cli.

```bash
drone secret add --image=plugins/slack \
    octocat/hello-world SLACK_WEBHOOK https://hooks.slack.com/services/...
```

Then sign the YAML file after all secrets are added.

```bash
drone sign octocat/hello-world
```

See [secrets](http://readme.drone.io/0.5/usage/secrets/) for additional
information on secrets

## Example

Common default example configuration:

```yaml
pipeline:
  slack:
    image: plugins/slack
    webhook: https://hooks.slack.com/services/...
    channel: dev
    username: drone
```

Example using a customized message:

```yaml
pipeline:
  slack:
    image: plugins/slack
    webhook: https://hooks.slack.com/services/...
    channel: dev
    username: drone
    template: >
      build #{{ build.number }} finished with a {{ build.status }} status
```

Example attach image inside a message:

```yaml
pipeline:
  slack:
    image: plugins/slack
    webhook: https://hooks.slack.com/services/...
    channel: dev
    username: drone
    template: >
      build #{{ build.number }} finished with a {{ build.status }} status
    image: https://cdn3.iconfinder.com/data/icons/picons-social/57/16-apple-128.png
```

Example change user avatar via icon URL:

```yaml
pipeline:
  slack:
    image: plugins/slack
    webhook: https://hooks.slack.com/services/...
    channel: dev
    username: drone
    template: >
      build #{{ build.number }} finished with a {{ build.status }} status
    icon_url: https://cdn0.iconfinder.com/data/icons/shift-free/32/Error-128.png
```

Example change user avatar via icon emoji:

```yaml
pipeline:
  slack:
    image: plugins/slack
    webhook: https://hooks.slack.com/services/...
    channel: dev
    username: drone
    template: >
      build #{{ build.number }} finished with a {{ build.status }} status
    icon_emoji: :+1:
```
