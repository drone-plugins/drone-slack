// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"testing"
	"time"

	"github.com/drone-plugins/drone-plugin-lib/drone"
	"gotest.tools/v3/assert"
)

func TestDefaultMessage(t *testing.T) {
	msg := defaultMessage(getTestPipeline())
	expectedMessage := "*success* <http://github.com/octocat/hello-world|octocat/hello-world#7fd1a60b> (master) by octocat"

	assert.Equal(t, expectedMessage, msg)
}

func TestShouldSendMessage(t *testing.T){
	assert.Equal(t, shouldSendMessage("does not match", "[A-Z].*"), false)
	assert.Equal(t, shouldSendMessage("Starts with capital letter - will match", "[A-Z].*"), true)
	assert.Equal(t, shouldSendMessage("The regexp is empty string - should still pass", ""), true)
}

func TestDefaultFallbackMessage(t *testing.T) {
	msg := fallback(getTestPipeline())
	expectedMessage := "success octocat/hello-world#7fd1a60b (master) by octocat"

	assert.Equal(t, expectedMessage, msg)
}

func TestTemplateMessage(t *testing.T) {
	s := getTestSettings()

	msg, err := templateMessage(s.Template, getTestPipeline())
	assert.NilError(t, err, "should create message by template without error")
	expectedMessage := `Message Template:
Initial commit

Message body
Initial commit
Message body`

	assert.Equal(t, expectedMessage, msg)
}

func TestTemplateFallbackMessage(t *testing.T) {
	s := getTestSettings()

	msg, err := templateMessage(s.Fallback, getTestPipeline())
	assert.NilError(t, err, "should create message by template without error")
	expectedMessage := `Message Template Fallback:
Initial commit
master
success`

	assert.Equal(t, expectedMessage, msg)
}

func getTestPipeline() drone.Pipeline {
	return drone.Pipeline{
		Repo:   getTestRepo(),
		Build:  getTestBuild(),
		Commit: getTestCommit(),
	}
}

func getTestRepo() drone.Repo {
	return drone.Repo{
		Owner: "octocat",
		Name:  "hello-world",
	}
}

func getTestBuild() drone.Build {
	return drone.Build{
		Tag:      "1.0.0",
		Event:    "push",
		Number:   1,
		Branch:   "master",
		DeployTo: "",
		Status:   "success",
		Link:     "http://github.com/octocat/hello-world",
		Started:  time.Unix(1546340400, 0), // 2019-01-01 12:00:00
		Created:  time.Unix(1546340400, 0), // 2019-01-01 12:00:00
	}
}

func getTestCommit() drone.Commit {
	return drone.Commit{
		SHA:     "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
		Ref:     "",
		Branch:  "master",
		Message: drone.ParseMessage("Initial commit\n\nMessage body"),
		Author: drone.Author{
			Username: "octocat",
			Name:     "The Octocat",
			Email:    "octocat@github.com",
			Avatar:   "https://avatars0.githubusercontent.com/u/583231?s=460&v=4",
		},
	}
}

func getTestSettings() Settings {
	t := `Message Template:
{{commit.message}}
{{commit.message.title}}
{{commit.message.body}}`

	tf := `Message Template Fallback:
{{commit.message.title}}
{{commit.branch}}
{{build.status}}`

	return Settings{
		Template: t,
		Fallback: tf,
	}
}
