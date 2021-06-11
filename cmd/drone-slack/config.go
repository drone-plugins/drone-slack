// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package main

import (
	"github.com/drone-plugins/drone-slack/plugin"
	"github.com/urfave/cli/v2"
)

// settingsFlags has the cli.Flags for the plugin.Settings.
func settingsFlags(settings *plugin.Settings) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "webhook",
			Usage:       "slack webhook url",
			EnvVars:     []string{"SLACK_WEBHOOK", "PLUGIN_WEBHOOK"},
			Destination: &settings.Webhook,
		},
		&cli.StringFlag{
			Name:        "channel",
			Usage:       "slack channel",
			EnvVars:     []string{"PLUGIN_CHANNEL"},
			Destination: &settings.Channel,
		},
		&cli.StringFlag{
			Name:        "recipient",
			Usage:       "slack recipient",
			EnvVars:     []string{"PLUGIN_RECIPIENT"},
			Destination: &settings.Recipient,
		},
		&cli.StringFlag{
			Name:        "username",
			Usage:       "slack username",
			EnvVars:     []string{"PLUGIN_USERNAME"},
			Destination: &settings.Username,
		},
		&cli.StringFlag{
			Name:        "template",
			Usage:       "slack template",
			EnvVars:     []string{"PLUGIN_TEMPLATE"},
			Destination: &settings.Template,
		},
		&cli.StringFlag{
			Name:        "only_allow_commits_regexp",
			Usage:       "regexp that should match commit message",
			EnvVars:     []string{"PLUGIN_TEMPLATE"},
			Destination: &settings.CommitAllowRegexp,
		},
		&cli.StringFlag{
			Name:        "fallback",
			Usage:       "slack fallback",
			EnvVars:     []string{"PLUGIN_FALLBACK"},
			Destination: &settings.Fallback,
		},
		&cli.StringFlag{
			Name:        "image",
			Usage:       "slack image url",
			EnvVars:     []string{"PLUGIN_IMAGE_URL"},
			Destination: &settings.ImageURL,
		},
		&cli.StringFlag{
			Name:        "icon.url",
			Usage:       "slack icon url",
			EnvVars:     []string{"PLUGIN_ICON_URL"},
			Destination: &settings.IconURL,
		},
		&cli.StringFlag{
			Name:        "icon.emoji",
			Usage:       "slack emoji url",
			EnvVars:     []string{"PLUGIN_ICON_EMOJI"},
			Destination: &settings.IconEmoji,
		},
		&cli.StringFlag{
			Name:        "color",
			Usage:       "slack color",
			EnvVars:     []string{"PLUGIN_COLOR"},
			Destination: &settings.Color,
		},
	}
}
