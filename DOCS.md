Use the Slack plugin to send a message to your Slack channel when a build completes.
You will need to supply Drone with an Incoming Webhook URL. You can create a new
Webhook URL here: https://my.slack.com/services/new/incoming-webhook

The following parameters are used to configuration the notification:

* **webhook_url** - json payloads are sent here
* **channel** - messages sent to the above webhook are posted here
* **recipient** - alternatively you can send it to a specific user
* **username** - choose the username this integration will post as
* **template** - Go template to create a custom message. See [docs](https://golang.org/pkg/text/template/)

The following is a sample Slack configuration in your .drone.yml file:

```yaml
notify:
  slack:
    webhook_url: https://hooks.slack.com/services/...
    channel: dev
    username: drone
    template: "*{{.Build.Status}}* {{.Repo.FullName}}#{{ShortCommit .Build}} ({{.Build.Branch}}) took {{Duration .Build}}"
```
