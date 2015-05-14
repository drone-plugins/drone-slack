To configure Slack Notifications, you first need to create a new Incoming WebHook.
Once you have the WebHook URL, you add the following to your drone.yml

```yaml
notify:
  slack:
    webhook_url: https://hooks.slack.com/services/...
    channel: #dev
    username: drone
```
