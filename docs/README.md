To configure Slack Notifications, you first need to create a new Incoming WebHook.
Once you have the WebHook URL, you add the following to your drone.yml

```yaml
notify:
  slack:
    webhook_url: https://hooks.slack.com/services/...
    channel: #dev
    username: drone
```

---

## Create a WebHook

In order to create a webhook you’ll need to login to your Slack account and navigate to: https://my.slack.com/services/new/incoming-webhook

Then follow the intructions to generate a new token:

![slack token screen](setup.png)

