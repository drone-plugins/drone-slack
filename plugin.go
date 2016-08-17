package main

import (
	"fmt"
	"math/rand"
	"strings"

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
		Webhook   string
		Channel   string
		Recipient string
		Username  string
		Template  string
		Success   MessageOptions
		Failure   MessageOptions
	}

	MessageOptions struct {
		Icon             string
		Username         string
		Template         string
		ImageAttachments []string
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
	}
)

func (p Plugin) Exec() error {
	var messageOptions MessageOptions

	attachment := slack.Attachment{
		Text:       message(p.Repo, p.Build),
		Fallback:   fallback(p.Repo, p.Build),
		Color:      color(p.Build),
		MarkdownIn: []string{"text", "fallback"},
	}

	payload := slack.WebHookPostPayload{}
	payload.Username = p.Config.Username
	payload.Attachments = []*slack.Attachment{&attachment}

	if p.Config.Recipient == "" {
		payload.Channel = prepend("#", p.Config.Channel)
	} else {
		payload.Channel = prepend("@", p.Config.Recipient)
	}

	if p.Config.Template != "" {
		txt, err := RenderTrim(p.Config.Template, p)
		if err != nil {
			return err
		}
		attachment.Text = txt
	}

	// Determine if the build was a success
	if p.Build.Status == "success" {
		messageOptions = p.Config.Success
	} else {
		messageOptions = p.Config.Failure
	}

	if strings.HasPrefix(messageOptions.Icon, "http") {
		payload.IconUrl = messageOptions.Icon
	} else {
		payload.IconEmoji = messageOptions.Icon
	}

	// Add image if any are provided
	imageCount := len(messageOptions.ImageAttachments)

	if imageCount > 0 {
		attachment.ImageURL = messageOptions.ImageAttachments[rand.Intn(imageCount)]
	}

	if messageOptions.Username != "" {
		payload.Username = messageOptions.Username
	}

	if messageOptions.Template != "" {
		txt, err := RenderTrim(messageOptions.Template, p)
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
