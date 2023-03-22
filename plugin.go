package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/drone/drone-template-lib/template"
	"github.com/slack-go/slack"
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
		Webhook   string
		Channel   string
		Mapping   string
		Username  string
		Template  string
		Fallback  string
		ImageURL  string
		IconURL   string
		IconEmoji string
		Color     string
		LinkNames bool
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
	attachment := slack.Attachment{
		Color:      p.Config.Color,
		ImageURL:   p.Config.ImageURL,
		AuthorName: "drone-slack",
		MarkdownIn: []string{"text", "fallback"},
	}
	if p.Config.Color == "" {
		attachment.Color = color(p.Build)
	}
	if p.Config.Fallback != "" {
		f, err := templateMessage(p.Config.Fallback, p)
		if err != nil {
			return err
		}
		attachment.Fallback = f
	} else {
		attachment.Fallback = fallback(p.Repo, p.Build)
	}

	if p.Config.Template != "" {
		var err error
		f, err := templateMessage(p.Config.Template, p)
		if err != nil {
			return err
		}
		attachment.Text = f
	} else {
		attachment.Text = message(p.Repo, p.Build)
	}

	payload := slack.WebhookMessage{
		Username:    p.Config.Username,
		Attachments: []slack.Attachment{attachment},
		IconURL:     p.Config.IconURL,
		IconEmoji:   p.Config.IconEmoji,
	}

	if p.Config.Mapping != "" {
		mapping, err := userMapping(p.Config.Mapping)
		if err != nil {
			return err
		}

		memberId, ok := mapping[p.Build.Author.Username]
		if !ok {
			return fmt.Errorf("user not found in mapping: %s", p.Build.Author.Username)
		}

		payload.Channel = memberId
	} else if p.Config.Channel != "" {
		payload.Channel = prepend("#", p.Config.Channel)
	}

	return slack.PostWebhook(p.Config.Webhook, &payload)
}

// userMapping gets the user mapping file.
func userMapping(value string) (map[string]string, error) {
	mapping := []byte(contents(value))

	// turn into a map
	values := map[string]string{}
	err := json.Unmarshal(mapping, &values)

	if err != nil {
		return nil, err
	}

	return values, nil
}

// contents gets the value referenced either in a local filem, a URL or the
// string value itself.
func contents(s string) string {
	if _, err := os.Stat(s); err == nil {
		o, _ := os.ReadFile(s)
		return os.ExpandEnv(string(o))
	}
	if _, err := url.Parse(s); err == nil {
		resp, err := http.Get(s)
		if err != nil {
			return s
		}
		defer resp.Body.Close()
		o, _ := io.ReadAll(resp.Body)
		return os.ExpandEnv(string(o))
	}
	return os.ExpandEnv(s)
}

func templateMessage(t string, plugin Plugin) (string, error) {
	return template.RenderTrim(t, plugin)
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
