package main

import (
	"fmt"
	"github.com/drone/drone-plugin-go/plugin"
)

type Slack struct {
	Webhook  string `json:"webhook_url"`
	Channel  string `json:"channel"`
	Username string `json:"username"`
}

func main() {
	repo := plugin.Repo{}
	links := plugin.Links{}
	commit := plugin.Commit{}
	vargs := Slack{}

	plugin.Param("links", &links)
	plugin.Param("commit", &commit)
	plugin.Param("repo", &repo)
	plugin.Param("vargs", &vargs)
	plugin.Parse()

	// create the Slack client
	client := Client{}
	client.Url = vargs.Webhook

	// generate the Slack message
	msg := Message{}
	msg.Channel = vargs.Channel
	msg.Username = vargs.Username

	attach := msg.NewAttachment()
	attach.Text = GetMessage(&links, &repo, &commit)
	attach.Fallback = GetFallback(&links, &repo, &commit)
	attach.Color = GetColor(&commit)
	attach.MrkdwnIn = []string{"text", "fallback"}

	// sends the message
	err := client.SendMessage(&msg)
	if err != nil {
		fmt.Println(err)
	}
}

func GetMessage(links *plugin.Links, repo *plugin.Repo, commit *plugin.Commit) string {
	return fmt.Sprintf("*%s* <%s|%s/%s#%s> (%s) by %s",
		commit.Status,
		links.Commit,
		repo.Owner,
		repo.Name,
		commit.Sha[:8],
		commit.Branch,
		commit.Author,
	)
}

func GetFallback(links *plugin.Links, repo *plugin.Repo, commit *plugin.Commit) string {
	return fmt.Sprintf("%s %s/%s#%s (%s) by %s",
		commit.Status,
		repo.Owner,
		repo.Name,
		commit.Sha[:8],
		commit.Branch,
		commit.Author,
	)
}

func GetColor(commit *plugin.Commit) string {
	switch commit.Status {
	case plugin.StatusSuccess:
		return "good"
	case plugin.StatusFailure, plugin.StatusError, plugin.StatusKilled:
		return "danger"
	default:
		return "warning"
	}
}
