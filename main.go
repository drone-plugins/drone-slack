package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/drone/drone-plugin-go/plugin"
)

type Slack struct {
	Webhook   string `json:"webhook_url"`
	Channel   string `json:"channel"`
	Recipient string `json:"recipient"`
	Username  string `json:"username"`
	Template  string `json:"template"`
}

// create the Slack client
var client ClientImpl = &Client{}

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

	client.SetUrl(vargs.Webhook)

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
	attach.Text = GetMessage(repo, build, system, vargs)
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

var defaultTemplate = "*{{.Build.Status}}* <{{.System.Link}}/{{.Repo.FullName}}/{{.Build.Number}}|{{.Repo.FullName}}#{{ShortCommit .Build}}> ({{.Build.Branch}}) by {{.Build.Author}}"

func GetMessage(repo plugin.Repo, build plugin.Build, system plugin.System, vargs Slack) string {
	templateString := defaultTemplate
	if len(vargs.Template) != 0 {
		templateString = vargs.Template
	}
	t := template.New("_").Funcs(template.FuncMap{
		"Duration": func(build plugin.Build) time.Duration {
			return time.Duration(build.Finished-build.Started) * time.Second
		},
		"ShortCommit": func(build plugin.Build) string { return build.Commit[:8] },
	})
	t, err := t.Parse(templateString)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	data := struct {
		Build  plugin.Build
		System plugin.System
		Repo   plugin.Repo
	}{build, system, repo}
	var textBuf bytes.Buffer
	if err := t.Execute(&textBuf, &data); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return textBuf.String()
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
