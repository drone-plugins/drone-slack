package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/drone/drone-plugin-go/plugin"
)

type Slack struct {
	Webhook   string `json:"webhook_url"`
	Channel   string `json:"channel"`
	Recipient string `json:"recipient"`
	Username  string `json:"username"`
}

func main() {
	repo := plugin.Repo{}
	build := plugin.Build{}
	system := plugin.System{}
	vargs := Slack{}

	plugin.Param("build", &build)
	plugin.Param("repo", &repo)
	plugin.Param("system", &system)
	plugin.Param("vargs", &vargs)

	// parse the parameters
	if err := plugin.Parse(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// create the Slack client
	client := Client{}
	client.Url = vargs.Webhook

	// generate the Slack message
	msg := Message{}
	msg.Username = vargs.Username
	msg.Channel = vargs.Recipient

	if len(vargs.Recipient) != 0 {
		msg.Channel = Prepend("@", vargs.Recipient)
	} else {
		msg.Channel = Prepend("#", vargs.Channel)
	}

	attach := msg.NewAttachment()
	attach.Text = GetMessage(&repo, &build, &system)
	attach.Fallback = GetFallback(&repo, &build)
	attach.Color = GetColor(&build)
	attach.MrkdwnIn = []string{"text", "fallback"}

	// sends the message
	if err := client.SendMessage(&msg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Prepend(prefix, s string) string {
	if !strings.HasPrefix(s, prefix) {
		return prefix + s
	}
	return s
}

func GetMessage(repo *plugin.Repo, build *plugin.Build, sys *plugin.System) string {
	return fmt.Sprintf("*%s* <%s|%s/%s#%s> (%s) by %s",
		build.Status,
		fmt.Sprintf("%s/%s/%v", sys.Link, repo.FullName, build.Number),
		repo.Owner,
		repo.Name,
		build.Commit[:8],
		build.Branch,
		build.Author,
	)
}

func GetFallback(repo *plugin.Repo, build *plugin.Build) string {
	return fmt.Sprintf("%s %s/%s#%s (%s) by %s",
		build.Status,
		repo.Owner,
		repo.Name,
		build.Commit[:8],
		build.Branch,
		build.Author,
	)
}

func GetColor(build *plugin.Build) string {
	switch build.Status {
	case plugin.StateSuccess:
		return "good"
	case plugin.StateFailure, plugin.StateError, plugin.StateKilled:
		return "danger"
	default:
		return "warning"
	}
}
