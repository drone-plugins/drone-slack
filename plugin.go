package main

import (
	"fmt"
	"strings"

	"github.com/bluele/slack"
	"github.com/drone/drone-template-lib/template"
)

type (
	Repo struct {
		Owner string
		Name  string
	}

	Build struct {
		Tag      string
		Event    string
		Number   int
		Commit   string
		Ref      string
		Branch   string
		Author   string
		Pull     string
		Message  string
		DeployTo string
		Status   string
		Link     string
		Started  int64
		Created  int64
	}

	Config struct {
		Webhook   string
		Channel   string
		Recipient string
		Username  string
		Template  string
		ImageURL  string
		IconURL   string
		IconEmoji string
		LinkNames bool
	}

	Job struct {
		Started int64
	}

	Computed struct {
		AuthorSlack    string
		RecipientSlack string
	}

	Plugin struct {
		Repo     Repo
		Build    Build
		Config   Config
		Job      Job
		Computed Computed
	}
)

func (p Plugin) Exec() error {
	attachment := slack.Attachment{
		Text:       message(p.Repo, p.Build),
		Fallback:   fallback(p.Repo, p.Build),
		Color:      color(p.Build),
		MarkdownIn: []string{"text", "fallback"},
		ImageURL:   p.Config.ImageURL,
	}

	payload := slack.WebHookPostPayload{}
	payload.Username = p.Config.Username
	payload.Attachments = []*slack.Attachment{&attachment}
	payload.IconUrl = p.Config.IconURL
	payload.IconEmoji = p.Config.IconEmoji

	if p.Computed.RecipientSlack != "" {
		payload.Channel = prepend("@", p.Computed.RecipientSlack)
	} else if p.Config.Channel != "" {
		payload.Channel = prepend("#", p.Config.Channel)
	}
	if p.Config.LinkNames == true {
		payload.LinkNames = "1"
	}
	if p.Config.Template != "" {
		txt, err := template.RenderTrim(p.Config.Template, p)

		if err != nil {
			return err
		}

		attachment.Text = txt
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

func prepend(prefix, s string) string {
	if !strings.HasPrefix(s, prefix) {
		return prefix + s
	}

	return s
}
