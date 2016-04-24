package main

import (
	"fmt"

	"github.com/bluele/slack"
)

type (
	Repo struct {
		Owner string
		Name  string
	}

	Build struct {
		Event  string
		Number int
		Commit string
		Branch string
		Author string
		Status string
		Link   string
	}

	Config struct {
		Webhook  string
		Channel  string
		Username string
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
	}
)

func (p Plugin) Exec() error {
	payload := slack.WebHookPostPayload{}
	payload.Channel = p.Config.Channel
	payload.Username = p.Config.Username
	payload.Attachments = []*slack.Attachment{
		{
			Text:       message(p.Repo, p.Build),
			Fallback:   fallback(p.Repo, p.Build),
			Color:      color(p.Build),
			MarkdownIn: []string{"text", "fallback"},
		},
	}

	client := slack.NewWebHook(p.Config.Webhook)
	return client.PostMessage(&payload)
}

func message(repo Repo, build Build) string {
	return fmt.Sprintf("*%s* <%s|%s/%s#%s> (%s) by %s",
		build.Status,
		build.Link,
		repo.Owner,
		repo.Name,
		build.Commit[:8],
		build.Branch,
		build.Author,
	)
}

func fallback(repo Repo, build Build) string {
	return fmt.Sprintf("%s %s/%s#%s (%s) by %s",
		build.Status,
		repo.Owner,
		repo.Name,
		build.Commit[:8],
		build.Branch,
		build.Author,
	)
}

func color(build Build) string {
	switch build.Status {
	case "success":
		return "good"
	case "failure", "error", "killed":
		return "danger"
	default:
		return "warning"
	}
}
