package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/drone/drone-template-lib/template"
	"github.com/slack-go/slack"
)

type (
	Repo struct {
		Owner string
		Name  string
	}

	BlockSet struct {
		Blocks []json.RawMessage `json:"blocks"`
	}

	Build struct {
		Tag      string
		Event    string
		Number   int
		Parent   int
		Commit   string
		Ref      string
		Branch   string
		Author   Author
		Pull     string
		Message  Message
		DeployTo string
		Status   string
		Link     string
		Started  int64
		Created  int64
	}

	Author struct {
		Username string
		Name     string
		Email    string
		Avatar   string
	}

	Message struct {
		msg   string
		Title string
		Body  string
	}

	Config struct {
		Webhook     string
		Channel     string
		Recipient   string
		Username    string
		Template    string
		Fallback    string
		ImageURL    string
		IconURL     string
		IconEmoji   string
		Color       string
		LinkNames   bool
		CustomBlock string
		AccessToken string
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

func (a Author) String() string {
	return a.Username
}

func newCommitMessage(m string) Message {
	// not checking the length here
	// as split will always return at least one element
	// check it if using more than the first element
	splitMsg := strings.Split(m, "\n")

	return Message{
		msg:   m,
		Title: strings.TrimSpace(splitMsg[0]),
		Body:  strings.TrimSpace(strings.Join(splitMsg[1:], "\n")),
	}
}

func (m Message) String() string {
	return m.msg
}

func (p Plugin) Exec() error {
	var blocks []slack.Block
	var channel string
	var text string
	var fallbackText string

	// Determine the channel
	if p.Config.Recipient != "" {
		channel = prepend("@", p.Config.Recipient)
	} else if p.Config.Channel != "" {
		channel = prepend("#", p.Config.Channel)
	}

	// Determine the message and fallback
	if p.Config.Template != "" {
		var err error
		text, err = templateMessage(p.Config.Template, p)
		if err != nil {
			return err
		}
	} else {
		text = message(p.Repo, p.Build)
	}

	if p.Config.Fallback != "" {
		var err error
		fallbackText, err = templateMessage(p.Config.Fallback, p)
		if err != nil {
			return err
		}
	} else {
		fallbackText = fallback(p.Repo, p.Build)
	}

	// Determine the color
	colorText := p.Config.Color
	if colorText == "" {
		colorText = color(p.Build)
	}

	// Parse custom blocks if they exist
	if p.Config.CustomBlock != "" {
		var blockSet BlockSet
		err := json.Unmarshal([]byte(p.Config.CustomBlock), &blockSet)
		if err != nil {
			return fmt.Errorf("could not unmarshal custom block: %w", err)
		}
		for _, rawMessage := range blockSet.Blocks {
			block := new(slack.SectionBlock)
			err := json.Unmarshal(rawMessage, block)
			if err != nil {
				return fmt.Errorf("could not unmarshal individual block: %w", err)
			}
			blocks = append(blocks, block)
		}
	}

	// If access token is provided, use it
	if p.Config.AccessToken != "" {
		slackApi := slack.New(p.Config.AccessToken)
		_, err := slackApi.AuthTest()
		if err != nil {
			return fmt.Errorf("failed to authenticate using access token: %w", err)
		}

		options := []slack.MsgOption{}
		if text != "" {
			options = append(options, slack.MsgOptionText(text, false))
		}
		if len(blocks) > 0 {
			options = append(options, slack.MsgOptionBlocks(blocks...))
		}

		_, _, err = slackApi.PostMessage(channel, options...)
		if err != nil {
			return fmt.Errorf("failed to post message using access token: %w", err)
		}
		return nil
	}

	// Build the attachment
	attachment := slack.Attachment{
		Color:      colorText,
		ImageURL:   p.Config.ImageURL,
		MarkdownIn: []string{"text", "fallback"},
		Text:       text,
		Fallback:   fallbackText,
	}

	// Build the payload
	payload := slack.WebhookMessage{
		Username:    p.Config.Username,
		Attachments: []slack.Attachment{attachment},
		IconURL:     p.Config.IconURL,
		IconEmoji:   p.Config.IconEmoji,
		Channel:     channel,
	}

	// Add custom blocks to the payload if they exist
	if len(blocks) > 0 {
		payload.Blocks = &slack.Blocks{
			BlockSet: blocks,
		}
	}

	// Post the message with the webhook
	return slack.PostWebhook(p.Config.Webhook, &payload)
}

func templateMessage(t string, plugin Plugin) (string, error) {
	c, err := contents(t)
	if err != nil {
		return "", fmt.Errorf("could not read template: %w", err)
	}

	return template.RenderTrim(c, plugin)
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
