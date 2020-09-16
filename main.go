package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli"
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
			EnvVar: "DRONE_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA",
			Value:  "00000000",
		},
		cli.StringFlag{
			Name:   "commit.ref",
			Value:  "refs/heads/master",
			Usage:  "git commit ref",
			EnvVar: "DRONE_COMMIT_REF",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "git commit branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
		cli.StringFlag{
			Name:   "commit.author",
			Usage:  "git author username",
			EnvVar: "DRONE_COMMIT_AUTHOR",
		},
		cli.StringFlag{
			Name:   "commit.author.email",
			Usage:  "git author email",
			EnvVar: "DRONE_COMMIT_AUTHOR_EMAIL",
		},
		cli.StringFlag{
			Name:   "commit.author.avatar",
			Usage:  "git author avatar",
			EnvVar: "DRONE_COMMIT_AUTHOR_AVATAR",
		},
		cli.StringFlag{
			Name:   "commit.author.name",
			Usage:  "git author name",
			EnvVar: "DRONE_COMMIT_AUTHOR_NAME",
		},
		cli.StringFlag{
			Name:   "commit.pull",
			Usage:  "git pull request",
			EnvVar: "DRONE_PULL_REQUEST",
		},
		cli.StringFlag{
			Name:   "commit.message",
			Usage:  "commit message",
			EnvVar: "DRONE_COMMIT_MESSAGE",
		},
		cli.StringFlag{
			Name:   "build.event",
			Value:  "push",
			Usage:  "build event",
			EnvVar: "DRONE_BUILD_EVENT",
		},
		cli.IntFlag{
			Name:   "build.number",
			Usage:  "build number",
			EnvVar: "DRONE_BUILD_NUMBER",
		},
		cli.StringFlag{
			Name:   "build.status",
			Usage:  "build status",
			Value:  "success",
			EnvVar: "DRONE_BUILD_STATUS",
		},
		cli.StringFlag{
			Name:   "build.link",
			Usage:  "build link",
			EnvVar: "DRONE_BUILD_LINK",
		},
		cli.Int64Flag{
			Name:   "build.started",
			Usage:  "build started",
			EnvVar: "DRONE_BUILD_STARTED",
		},
		cli.Int64Flag{
			Name:   "build.created",
			Usage:  "build created",
			EnvVar: "DRONE_BUILD_CREATED",
		},
		cli.StringFlag{
			Name:   "build.tag",
			Usage:  "build tag",
			EnvVar: "DRONE_TAG",
		},
		cli.StringFlag{
			Name:   "build.deployTo",
			Usage:  "environment deployed to",
			EnvVar: "DRONE_DEPLOY_TO",
		},
		cli.Int64Flag{
			Name:   "job.started",
			Usage:  "job started",
			EnvVar: "DRONE_JOB_STARTED",
		},

		cli.StringFlag{
			Name:   "build.failedStages",
			Usage:  "comma separated list of failed stages",
			EnvVar: "DRONE_FAILED_STAGES",
		},
		cli.StringFlag{
			Name:   "build.failedSteps",
			Usage:  "comma separated list of failed steps",
			EnvVar: "DRONE_FAILED_STEPS",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "enable debug logs",
			EnvVar: "DEBUG",
		},
	}

	if _, err := os.Stat("/run/drone/env"); err == nil {
		godotenv.Overload("/run/drone/env")
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
			Tag:          c.String("build.tag"),
			Number:       c.Int("build.number"),
			Event:        c.String("build.event"),
			Status:       c.String("build.status"),
			FailedSteps:  c.String("build.failedSteps"),
			FailedStages: c.String("build.failedStages"),
			Commit:       c.String("commit.sha"),
			Ref:          c.String("commit.ref"),
			Branch:       c.String("commit.branch"),
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
			Webhook:   c.String("webhook"),
			Channel:   c.String("channel"),
			Recipient: c.String("recipient"),
			Username:  c.String("username"),
			Template:  c.String("template"),
			Fallback:  c.String("fallback"),
			ImageURL:  c.String("image"),
			IconURL:   c.String("icon.url"),
			IconEmoji: c.String("icon.emoji"),
			Color:     c.String("color"),
			LinkNames: c.Bool("link-names"),
			Debug:     c.Bool("debug"),
		},
	}

	if plugin.Config.Webhook == "" {
		return errors.New("Missing webhook")
	}

	return plugin.Exec()
}
