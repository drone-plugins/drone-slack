package main

import "testing"
import "gotest.tools/assert"

func TestDefaultMessage(t *testing.T) {
	repo := getTestRepo()
	build := getTestBuild()

	msg := message(repo, build)
	expectedMessage := "*success* <http://github.com/octocat/hello-world|octocat/hello-world#7fd1a60b> (master) by octocat"

	assert.Equal(t, expectedMessage, msg)
}

func TestFallbackMessage(t *testing.T) {
	repo := getTestRepo()
	build := getTestBuild()

	msg := fallback(repo, build)
	expectedMessage := "success octocat/hello-world#7fd1a60b (master) by octocat"

	assert.Equal(t, expectedMessage, msg)
}

func getTestRepo() Repo {
	return Repo{
		Owner: "octocat",
		Name:  "hello-world",
	}
}

func getTestBuild() Build {
	author := Author{
		Username: "octocat",
		Name:     "The Octocat",
		Email:    "octocat@github.com",
		Avatar:   "https://avatars0.githubusercontent.com/u/583231?s=460&v=4",
	}

	return Build{
		Tag:      "1.0.0",
		Event:    "push",
		Number:   1,
		Commit:   "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
		Ref:      "",
		Branch:   "master",
		Author:   author,
		Pull:     "",
		Message:  Message{
			msg: "Initial commit",
			Title: "Initial commit",
		},
		DeployTo: "",
		Status:   "success",
		Link:     "http://github.com/octocat/hello-world",
		Started:  1546340400, // 2019-01-01 12:00:00
		Created:  1546340400, // 2019-01-01 12:00:00
	}
}
