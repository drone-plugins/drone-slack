package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/drone-plugins/drone-slack/slack"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin"
	"github.com/drone/drone-go/template"
)

type Slack struct {
	Webhook   string `json:"webhook_url"`
	Channel   string `json:"channel"`
	Recipient string `json:"recipient"`
	Username  string `json:"username"`
	Template  string `json:"template"`
}

func main() {
	var (
		repo  = new(drone.Repo)
		build = new(drone.Build)
		sys   = new(drone.System)
		vargs = new(Slack)
	)

	plugin.Param("build", build)
	plugin.Param("repo", repo)
	plugin.Param("system", sys)
	plugin.Param("vargs", vargs)

	err := plugin.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// create the Slack client
	client := slack.NewClient(vargs.Webhook)

	// generate the Slack message
	msg := slack.Message{
		Username: vargs.Username,
		Channel:  vargs.Recipient,
	}

	// prepend the @ or # symbol if the user forgot to include
	// in their configuration string.
	if len(vargs.Recipient) != 0 {
		msg.Channel = prepend("@", vargs.Recipient)
	} else {
		msg.Channel = prepend("#", vargs.Channel)
	}

	attach := msg.NewAttachment()
	attach.Text = message(repo, build, sys)
	attach.Fallback = fallback(repo, build)
	attach.Color = color(build)
	attach.MrkdwnIn = []string{"text", "fallback"}

	// the user may choose to override the template, if so parse
	// and execute the template to override our default message
	if len(vargs.Template) != 0 {
		attach.Text, err = template.RenderTrim(
			vargs.Template,
			&drone.Payload{
				Build:  build,
				Repo:   repo,
				System: sys,
			},
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// sends the message
	if err := client.SendMessage(&msg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func prepend(prefix, s string) string {
	if !strings.HasPrefix(s, prefix) {
		return prefix + s
	}
	return s
}

func message(repo *drone.Repo, build *drone.Build, sys *drone.System) string {
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

func fallback(repo *drone.Repo, build *drone.Build) string {
	return fmt.Sprintf("%s %s/%s#%s (%s) by %s",
		build.Status,
		repo.Owner,
		repo.Name,
		build.Commit[:8],
		build.Branch,
		build.Author,
	)
}

func color(build *drone.Build) string {
	switch build.Status {
	case drone.StatusSuccess:
		return "good"
	case drone.StatusFailure, drone.StatusError, drone.StatusKilled:
		return "danger"
	default:
		return "warning"
	}
}
