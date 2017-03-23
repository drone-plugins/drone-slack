package main

import (
	"fmt"
	"strings"

	"github.com/bluele/slack"
)

type (
	Repo struct {
		Owner string
		Name  string
	}

	Build struct {
		Tag        string
		Event      string
		Number     int
		Commit     string
		Ref        string
		Branch     string
		Author     string
		Message    string
		CommitLink string
		Status     string
		Link       string
		Started    int64
		Created    int64
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
	}

	Job struct {
		Started int64
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
		Job    Job
	}
)

func (p Plugin) Exec() error {
	attachment := slack.Attachment{
		Fallback:   fallback(p.Repo, p.Build),
		Color:      color(p.Build),
		MarkdownIn: []string{"text", "fallback"},
		ImageURL:   p.Config.ImageURL,
		Fields:     fields(p.Repo, p.Build),
	}

	payload := slack.WebHookPostPayload{}
	payload.Username = p.Config.Username
	payload.Attachments = []*slack.Attachment{&attachment}
	payload.IconUrl = p.Config.IconURL
	payload.IconEmoji = p.Config.IconEmoji

	if p.Config.Recipient != "" {
		payload.Channel = prepend("@", p.Config.Recipient)
	} else if p.Config.Channel != "" {
		payload.Channel = prepend("#", p.Config.Channel)
	}

	if p.Config.Template != "" {
		txt, err := RenderTrim(p.Config.Template, p)
		if err != nil {
			return err
		}
		attachment.Text = txt
	}

	client := slack.NewWebHook(p.Config.Webhook)
	return client.PostMessage(&payload)
}

func fields(repo Repo, build Build) []*slack.AttachmentField {
	return []*slack.AttachmentField{
		&slack.AttachmentField{
			Title: "Status",
			Value: build.Status,
			Short: true,
		},
		&slack.AttachmentField{
			Title: "Branch",
			Value: build.Branch,
			Short: true,
		},
		&slack.AttachmentField{
			Title: "Output",
			Value: fmt.Sprintf("<%s|%s/%s#%d>", build.Link, repo.Owner, repo.Name, build.Number),
			Short: true,
		},
		&slack.AttachmentField{
			Title: "Commit",
			Value: fmt.Sprintf("<%s|%s>", build.CommitLink, build.Commit[:12]),
			Short: true,
		},
		&slack.AttachmentField{
			Title: "Author",
			Value: build.Author,
			Short: false,
		},
		&slack.AttachmentField{
			Title: "Commit Message",
			Value: build.Message,
			Short: false,
		},
	}
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
