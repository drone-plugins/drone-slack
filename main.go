package main

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
	"log"
	"os"
)

var (
	version = "0.0.0"
	build   = "0"
)

func main() {

	app := cli.NewApp()
	app.Name = "slack plugin"
	app.Usage = "slack plugin"
	app.Action = run
	app.Version = fmt.Sprintf("%s+%s", version, build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "webhook",
			Usage:  "slack webhook url",
			EnvVar: "SLACK_WEBHOOK,PLUGIN_WEBHOOK",
		},
		cli.StringFlag{
			Name:   "channel",
			Usage:  "slack channel",
			EnvVar: "PLUGIN_CHANNEL",
		},
		cli.StringFlag{
			Name:   "recipient",
			Usage:  "slack recipient",
			EnvVar: "PLUGIN_RECIPIENT",
		},
		cli.StringFlag{
			Name:   "username",
			Usage:  "slack username",
			EnvVar: "PLUGIN_USERNAME",
		},
		cli.StringFlag{
			Name:   "template",
			Usage:  "slack template",
			EnvVar: "PLUGIN_TEMPLATE",
		},
		cli.StringFlag{
			Name:   "fallback",
			Usage:  "slack fallback",
			EnvVar: "PLUGIN_FALLBACK",
		},
		cli.BoolFlag{
			Name:   "link-names",
			Usage:  "slack link names",
			EnvVar: "PLUGIN_LINK_NAMES",
		},
		cli.StringFlag{
			Name:   "image",
			Usage:  "slack image url",
			EnvVar: "PLUGIN_IMAGE_URL",
		},
		cli.StringFlag{
			Name:   "color",
			Usage:  "slack color",
			EnvVar: "PLUGIN_COLOR",
		},
		cli.StringFlag{
			Name:   "icon.url",
			Usage:  "slack icon url",
			EnvVar: "PLUGIN_ICON_URL",
		},
		cli.StringFlag{
			Name:   "icon.emoji",
			Usage:  "slack emoji url",
			EnvVar: "PLUGIN_ICON_EMOJI",
		},
		cli.StringFlag{
			Name:   "repo.owner",
			Usage:  "repository owner",
			EnvVar: "CI_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "CI_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "CI_COMMIT_SHA",
			Value:  "00000000",
		},
		cli.StringFlag{
			Name:   "commit.ref",
			Value:  "refs/heads/master",
			Usage:  "git commit ref",
			EnvVar: "CI_COMMIT_REF",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "git commit branch",
			EnvVar: "CI_COMMIT_BRANCH",
		},
		cli.StringFlag{
			Name:   "commit.author",
			Usage:  "git author username",
			EnvVar: "CI_COMMIT_AUTHOR",
		},
		cli.StringFlag{
			Name:   "commit.author.email",
			Usage:  "git author email",
			EnvVar: "CI_COMMIT_AUTHOR_EMAIL",
		},
		cli.StringFlag{
			Name:   "commit.author.avatar",
			Usage:  "git author avatar",
			EnvVar: "CI_COMMIT_AUTHOR_AVATAR",
		},
		cli.StringFlag{
			Name:   "commit.author.name",
			Usage:  "git author name",
			EnvVar: "CI_COMMIT_AUTHOR_NAME",
		},
		cli.StringFlag{
			Name:   "commit.pull",
			Usage:  "git pull request",
			EnvVar: "CI_COMMIT_PULL_REQUEST",
		},
		cli.StringFlag{
			Name:   "commit.message",
			Usage:  "commit message",
			EnvVar: "CI_COMMIT_MESSAGE",
		},
		cli.StringFlag{
			Name:   "build.event",
			Value:  "push",
			Usage:  "build event",
			EnvVar: "CI_PIPELINE_EVENT",
		},
		cli.IntFlag{
			Name:   "build.number",
			Usage:  "build number",
			EnvVar: "CI_PIPELINE_NUMBER",
		},
		cli.IntFlag{
			Name:   "build.parent",
			Usage:  "build parent",
			EnvVar: "CI_PIPELINE_PARENT",
		},
		cli.StringFlag{
			Name:   "build.status",
			Usage:  "build status",
			Value:  "success",
			EnvVar: "PLUGIN_STATUS",
		},
		cli.StringFlag{
			Name:   "build.link",
			Usage:  "build link",
			EnvVar: "CI_PIPELINE_LINK",
		},
		cli.Int64Flag{
			Name:   "build.started",
			Usage:  "build started",
			EnvVar: "CI_PIPELINE_STARTED",
		},
		cli.Int64Flag{
			Name:   "build.created",
			Usage:  "build created",
			EnvVar: "CI_PIPELINE_CREATED",
		},
		cli.StringFlag{
			Name:   "build.tag",
			Usage:  "build tag",
			EnvVar: "CI_COMMIT_TAG",
		},
		cli.StringFlag{
			Name:   "build.deployTo",
			Usage:  "environment deployed to",
			EnvVar: "CI_PIPELINE_DEPLOY_TARGET",
		},
		cli.Int64Flag{
			Name:   "job.started",
			Usage:  "job started",
			EnvVar: "CI_PIPELINE_STARTED",
		},
		cli.StringFlag{
			Name:   "custom.block",
			Usage:  "custom block to send to slack. ",
			EnvVar: "PLUGIN_CUSTOM_BLOCK",
		},
		cli.StringFlag{
			Name:   "access.token",
			Usage:  "slack access token",
			EnvVar: "PLUGIN_ACCESS_TOKEN,SLACK_ACCESS_TOKEN",
		},
		cli.StringFlag{
			Name:   "mentions",
			Usage:  "slack mentions for the message.",
			EnvVar: "PLUGIN_MENTIONS",
		},
		cli.StringFlag{
			Name:   "custom.template",
			Usage:  "prebuilt custom template for the message.",
			EnvVar: "PLUGIN_CUSTOM_TEMPLATE",
		},
		cli.StringFlag{
			Name:   "message",
			Usage:  "slack message. either this or the custom template must be set. ",
			EnvVar: "PLUGIN_MESSAGE",
		},

		// File send params
		cli.StringFlag{
			Name:   "filepath",
			Usage:  "slack file path",
			EnvVar: "PLUGIN_FILE_PATH",
		},
		cli.StringFlag{
			Name:   "filename",
			Usage:  "slack file name",
			EnvVar: "PLUGIN_FILE_NAME",
		},
		cli.StringFlag{
			Name:   "title",
			Usage:  "slack title",
			EnvVar: "PLUGIN_TITLE",
		},
		cli.StringFlag{
			Name:   "initial_comment",
			Usage:  "slack initial comment",
			EnvVar: "PLUGIN_INITIAL_COMMENT",
		},
		cli.BoolFlag{
			Name:   "fail_on_error",
			Usage:  "fail build on error",
			EnvVar: "PLUGIN_FAIL_ON_ERROR",
		},
		cli.StringFlag{
			Name:   "slack_id_of",
			Usage:  "slack id required for the user email id",
			EnvVar: "PLUGIN_SLACK_USER_EMAIL_ID",
		},
		cli.StringFlag{
			Name:   "committer_list_git_path",
			Usage:  "git repo path holding the committers email id to fetch slack IDs from",
			EnvVar: "PLUGIN_GIT_REPO_PATH",
		},
		cli.BoolFlag{
			Name:   "plugin_committer_slack_id",
			Usage:  "flag to enable fetching slack IDs from the committers list",
			EnvVar: "PLUGIN_COMMITTERS_SLACK_ID",
		},
	}

	if _, err := os.Stat("/run/drone/env"); err == nil {
		_ = godotenv.Overload("/run/drone/env")
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Repo: Repo{
			Owner: c.String("repo.owner"),
			Name:  c.String("repo.name"),
		},
		Build: Build{
			Tag:    c.String("build.tag"),
			Number: c.Int("build.number"),
			Parent: c.Int("build.parent"),
			Event:  c.String("build.event"),
			Status: c.String("build.status"),
			Commit: c.String("commit.sha"),
			Ref:    c.String("commit.ref"),
			Branch: c.String("commit.branch"),
			Author: Author{
				Username: c.String("commit.author"),
				Name:     c.String("commit.author.name"),
				Email:    c.String("commit.author.email"),
				Avatar:   c.String("commit.author.avatar"),
			},
			Pull:     c.String("commit.pull"),
			Message:  newCommitMessage(c.String("commit.message")),
			DeployTo: c.String("build.deployTo"),
			Link:     c.String("build.link"),
			Started:  c.Int64("build.started"),
			Created:  c.Int64("build.created"),
		},
		Job: Job{
			Started: c.Int64("job.started"),
		},
		Config: Config{
			Webhook:        c.String("webhook"),
			Channel:        c.String("channel"),
			Recipient:      c.String("recipient"),
			Username:       c.String("username"),
			Template:       c.String("template"),
			Fallback:       c.String("fallback"),
			ImageURL:       c.String("image"),
			IconURL:        c.String("icon.url"),
			IconEmoji:      c.String("icon.emoji"),
			Color:          c.String("color"),
			LinkNames:      c.Bool("link-names"),
			CustomBlock:    c.String("custom.block"),
			AccessToken:    c.String("access.token"),
			Mentions:       c.String("mentions"),
			CustomTemplate: c.String("custom.template"),
			Message:        c.String("message"),
			// File upload attributes
			FilePath:             c.String("filepath"),
			FileName:             c.String("filename"),
			Title:                c.String("title"),
			InitialComment:       c.String("initial_comment"),
			FailOnError:          c.Bool("fail_on_error"),
			SlackIdOf:            c.String("slack_id_of"),
			CommitterListGitPath: c.String("committer_list_git_path"),
			CommitterSlackId:     c.Bool("plugin_committer_slack_id"),
		},
	}

	if plugin.Build.Commit == "" {
		plugin.Build.Commit = "0000000000000000000000000000000000000000"
	}
	if plugin.Config.Webhook == "" && plugin.Config.AccessToken == "" {
		return errors.New("you must provide a webhook url or access token")
	}

	return plugin.Exec()
}
