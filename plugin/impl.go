// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/drone-plugins/drone-plugin-lib/drone"
	"github.com/drone/drone-template-lib/template"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type (
	// Settings for the plugin.
	Settings struct {
		Webhook           string
		Channel           string
		Recipient         string
		Username          string
		Template          string
		Fallback          string
		ImageURL          string
		IconURL           string
		IconEmoji         string
		Color             string
		CommitAllowRegexp string
	}
)

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	if p.settings.Webhook == "" {
		return errors.New("missing webhook")
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	if !shouldSendMessage(p.pipeline.Commit.Message.String(), p.settings.CommitAllowRegexp) {
		return nil
	}

	attachment := slack.Attachment{
		Color:      p.settings.Color,
		ImageURL:   p.settings.ImageURL,
		MarkdownIn: []string{"text", "fallback"},
	}
	if p.settings.Color == "" {
		attachment.Color = color(p.pipeline.Build)
	}

	if p.settings.Fallback != "" {
		f, err := templateMessage(p.settings.Fallback, p.pipeline)
		if err != nil {
			return fmt.Errorf("could not create fallback message: %w", err)
		}
		attachment.Fallback = f
	} else {
		attachment.Fallback = fallback(p.pipeline)
	}

	msg := slack.WebhookMessage{
		Username:    p.settings.Username,
		Attachments: []slack.Attachment{attachment},
		IconURL:     p.settings.IconURL,
		IconEmoji:   p.settings.IconEmoji,
	}

	if p.settings.Recipient != "" {
		msg.Channel = prepend("@", p.settings.Recipient)
	} else if p.settings.Channel != "" {
		msg.Channel = prepend("#", p.settings.Channel)
	}

	if p.settings.Template != "" {
		var err error
		attachment.Text, err = templateMessage(p.settings.Template, p.pipeline)
		if err != nil {
			return fmt.Errorf("could not create template message: %w", err)
		}
	} else {
		attachment.Text = defaultMessage(p.pipeline)
	}

	logrus.WithFields(logrus.Fields{
		"channel":  msg.Channel,
		"username": msg.Username,
	}).Info("sending message")
	err := slack.PostWebhookCustomHTTPContext(p.network.Context, p.settings.Webhook, p.network.Client, &msg)
	if err != nil {
		return fmt.Errorf("could not send webhook: %w", err)
	}

	return nil
}

func detectRef(build drone.Build, commit drone.Commit) string {
	if commit.SHA != "" {
		return commit.SHA[:8]
	}

	return build.Tag
}

func templateMessage(t string, p drone.Pipeline) (string, error) {
	return template.RenderTrim(t, p)
}


func shouldSendMessage(msg string, allowRegexp string) bool{
	if allowRegexp == "" {
		return true
	}else{
		var mustMatch = regexp.MustCompile(allowRegexp)
		return mustMatch.MatchString(msg)
	}
}

func defaultMessage(p drone.Pipeline) string {
	return fmt.Sprintf("*%s* <%s|%s/%s#%s> (%s) by %s",
		p.Build.Status,
		p.Build.Link,
		p.Repo.Owner,
		p.Repo.Name,
		detectRef(p.Build, p.Commit),
		p.Build.Branch,
		p.Commit.Author,
	)
}

func fallback(p drone.Pipeline) string {
	return fmt.Sprintf("%s %s/%s#%s (%s) by %s",
		p.Build.Status,
		p.Repo.Owner,
		p.Repo.Name,
		detectRef(p.Build, p.Commit),
		p.Build.Branch,
		p.Commit.Author,
	)
}

func color(build drone.Build) string {
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
