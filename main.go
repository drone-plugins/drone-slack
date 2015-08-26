package main

import (
	"fmt"
	"os"

	"github.com/drone/drone-plugin-go/plugin"
)

type Slack struct {
	Webhook  string `json:"webhook_url"`
	Channel  string `json:"channel"`
	Username string `json:"username"`
}

func main() {
	repo := plugin.Repo{}
	build := plugin.Build{}
	vargs := Slack{}

	plugin.Param("build", &build)
	plugin.Param("repo", &repo)
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
	msg.Channel = vargs.Channel
	msg.Username = vargs.Username

	attach := msg.NewAttachment()
	attach.Text = GetMessage(&repo, &build)
	attach.Fallback = GetFallback(&repo, &build)
	attach.Color = GetColor(&build)
	attach.MrkdwnIn = []string{"text", "fallback"}

	// sends the message
	if err := client.SendMessage(&msg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GetMessage(repo *plugin.Repo, build *plugin.Build) string {
	return fmt.Sprintf("*%s* <%s|%s/%s#%s> (%s) by %s",
		build.Status,
		fmt.Sprintf("%s/%v", repo.Self, build.Number),
		repo.Owner,
		repo.Name,
		build.Commit.Sha[:8],
		build.Commit.Branch,
		build.Commit.Author.Login,
	)
}

func GetFallback(repo *plugin.Repo, build *plugin.Build) string {
	return fmt.Sprintf("%s %s/%s#%s (%s) by %s",
		build.Status,
		repo.Owner,
		repo.Name,
		build.Commit.Sha[:8],
		build.Commit.Branch,
		build.Commit.Author.Login,
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
