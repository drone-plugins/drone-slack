package main

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestExec(t *testing.T) {
	plugin := Plugin{
		Repo:   getTestRepo(),
		Build:  getTestBuild(),
		Job:    getTestJob(),
		Config: getTestConfig(),
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		out, _ := io.ReadAll(r.Body)
		got := string(out)
		want := `{"attachments":[{"color":"good","fallback":"Message Template Fallback:\nInitial commit\nmaster\nsuccess","text":"Message Template:\nInitial commit\n\nMessage body\nInitial commit\nMessage body","mrkdwn_in":["text","fallback"],"blocks":null}],"replace_original":false,"delete_original":false}`
		assert.Equal(t, got, want)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	plugin.Config.Webhook = server.URL
	_ = plugin.Exec()
}

func TestNewCommitMessage(t *testing.T) {
	testCases := map[string]struct {
		Msg    string
		Expect Message
	}{
		"Empty Message": {
			Msg:    "",
			Expect: Message{},
		},
		"Title Only": {
			Msg: "title with space",
			Expect: Message{
				msg:   "title with space",
				Title: "title with space",
				Body:  "",
			},
		},
		"Title and Body": {
			Msg: "Title with space\nBody with space",
			Expect: Message{
				msg:   "Title with space\nBody with space",
				Title: "Title with space",
				Body:  "Body with space",
			},
		},
		"Empty Second Line": {
			Msg: "Title with space\n\nBody with space",
			Expect: Message{
				msg:   "Title with space\n\nBody with space",
				Title: "Title with space",
				Body:  "Body with space",
			},
		},
	}

	for name, testCase := range testCases {
		assert.Equal(t, testCase.Expect, newCommitMessage(testCase.Msg), name)
	}
}

func TestDefaultMessage(t *testing.T) {
	repo := getTestRepo()
	build := getTestBuild()

	msg := message(repo, build)
	expectedMessage := "*success* <http://github.com/octocat/hello-world|octocat/hello-world#7fd1a60b> (master) by octocat"

	assert.Equal(t, expectedMessage, msg)
}

func TestDefaultFallbackMessage(t *testing.T) {
	repo := getTestRepo()
	build := getTestBuild()

	msg := fallback(repo, build)
	expectedMessage := "success octocat/hello-world#7fd1a60b (master) by octocat"

	assert.Equal(t, expectedMessage, msg)
}

func TestTemplateMessage(t *testing.T) {
	plugin := getTestPlugin()

	msg, err := templateMessage(plugin.Config.Template, plugin)
	assert.NilError(t, err, "should create message by template without error")
	expectedMessage := `Message Template:
Initial commit

Message body
Initial commit
Message body`

	assert.Equal(t, expectedMessage, msg)
}

func TestTemplateFallbackMessage(t *testing.T) {
	plugin := getTestPlugin()

	msg, err := templateMessage(plugin.Config.Fallback, plugin)
	assert.NilError(t, err, "should create message by template without error")
	expectedMessage := `Message Template Fallback:
Initial commit
master
success`

	assert.Equal(t, expectedMessage, msg)
}

func getTestPlugin() Plugin {
	return Plugin{
		Repo:   getTestRepo(),
		Build:  getTestBuild(),
		Config: getTestConfig(),
	}
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
		Message:  newCommitMessage("Initial commit\n\nMessage body"),
		DeployTo: "",
		Status:   "success",
		Link:     "http://github.com/octocat/hello-world",
		Started:  1546340400, // 2019-01-01 12:00:00
		Created:  1546340400, // 2019-01-01 12:00:00
	}
}

func getTestJob() Job {
	return Job{
		Started: 1546340400,
	}
}

func getTestConfig() Config {
	t := `Message Template:
{{build.message}}
{{build.message.title}}
{{build.message.body}}`

	tf := `Message Template Fallback:
{{build.message.title}}
{{build.branch}}
{{build.status}}`

	return Config{
		Template: t,
		Fallback: tf,
	}
}

func TestFileUpload(t *testing.T) {

	cfg := Config{
		Channel:        os.Getenv("PLUGIN_CHANNEL"),
		AccessToken:    os.Getenv("PLUGIN_ACCESS_TOKEN"),
		FilePath:       os.Getenv("PLUGIN_FILE_PATH"),
		FileName:       os.Getenv("PLUGIN_FILE_NAME"),
		InitialComment: os.Getenv("PLUGIN_INITIAL_COMMENT"),
		Title:          os.Getenv("PLUGIN_TITLE"),
	}

	plugin := Plugin{
		Repo:   getTestRepo(),
		Build:  getTestBuild(),
		Job:    getTestJob(),
		Config: cfg,
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		out, _ := io.ReadAll(r.Body)
		got := string(out)

		want := `{"attachments":[{"color":"good","fallback":"success octocat/hello-world#7fd1a60b (master) by octocat","text":"*success* \u003chttp://github.com/octocat/hello-world|octocat/hello-world#7fd1a60b\u003e (master) by octocat","mrkdwn_in":["text","fallback"],"blocks":null}],"replace_original":false,"delete_original":false}`
		assert.Equal(t, got, want)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	plugin.Config.Webhook = server.URL

	_ = plugin.Exec()
}

func TestGetSlackIdFromEmail(t *testing.T) {
	config := Config{
		AccessToken: "test-access-token",
		SlackIdOf:   "octocat@github.com",
	}

	plugin := Plugin{
		Config: config,
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := `{
			"user": {
				"id": "U12345",
				"name": "octocat"
			}
		}`
		_, _ = w.Write([]byte(response))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	plugin.Config.AccessToken = server.URL

	err := GetSlackIdFromEmail(&plugin)
	assert.NilError(t, err, "should retrieve Slack ID without error")
}

func TestGetSlackIdsOfCommitters(t *testing.T) {
	config := Config{
		AccessToken:          "mock-access-token",
		CommitterListGitPath: "/mock/repo/path",
	}

	plugin := Plugin{Config: config}

	mockGetAuthorsList := func(gitDir string) ([]string, error) {
		return []string{"user1@example.com", "user2@example.com"}, nil
	}

	mockGetSlackUserIDByEmail := func(accessToken string, emailList string) ([]string, error) {
		emails := strings.Split(emailList, ",")
		emailToID := map[string]string{
			"user1@example.com": "U12345",
			"user2@example.com": "U67890",
		}

		var ids []string
		for _, email := range emails {
			id, ok := emailToID[email]
			if !ok {
				return nil, fmt.Errorf("invalid_auth")
			}
			ids = append(ids, id)
		}
		return ids, nil
	}

	slackIDs, err := GetSlackIdsOfCommitters(&plugin, mockGetAuthorsList, mockGetSlackUserIDByEmail)
	assert.NilError(t, err, "should retrieve Slack IDs without error")
	assert.DeepEqual(t, slackIDs, []string{"U12345", "U67890"})
}

func TestGetSlackIdsOfCommitters_NoCommitters(t *testing.T) {
	config := Config{
		AccessToken:          "mock-access-token",
		CommitterListGitPath: "/mock/repo/path",
	}

	plugin := Plugin{Config: config}

	mockGetAuthorsList := func(gitDir string) ([]string, error) {
		return []string{}, nil
	}

	mockGetSlackUserIDByEmail := func(accessToken string, emailList string) ([]string, error) {
		if emailList == "" {
			return []string{}, nil
		}
		return nil, fmt.Errorf("unexpected call to getSlackUserIDByEmail")
	}

	slackIDs, err := GetSlackIdsOfCommitters(&plugin, mockGetAuthorsList, mockGetSlackUserIDByEmail)

	assert.NilError(t, err, "should handle a new repository without error")

	expectedSlackIDs := []string{}
	if diff := cmp.Diff(expectedSlackIDs, slackIDs); diff != "" {
		t.Errorf("mismatch in Slack IDs (-want +got):\n%s", diff)
	}
}

func TestGetSlackIdsOfCommitters_SlackUserLookupFailure(t *testing.T) {
	config := Config{
		AccessToken:          "mock-access-token",
		CommitterListGitPath: "/mock/repo/path",
	}

	plugin := Plugin{Config: config}
	mockGetAuthorsList := func(gitDir string) ([]string, error) {
		return []string{"user1@example.com", "user2@example.com"}, nil
	}

	mockGetSlackUserIDByEmail := func(accessToken string, emailList string) ([]string, error) {
		emails := strings.Split(emailList, ",")
		emailToID := map[string]string{
			"user1@example.com": "U12345",
		}

		var ids []string
		for _, email := range emails {
			id, ok := emailToID[email]
			if !ok {
				return nil, fmt.Errorf("user lookup failed for email: %s", email)
			}
			ids = append(ids, id)
		}
		return ids, nil
	}

	slackIDs, err := GetSlackIdsOfCommitters(&plugin, mockGetAuthorsList, mockGetSlackUserIDByEmail)
	assert.ErrorContains(t, err, "user lookup failed for email: user2@example.com")
	expectedSlackIDs := []string{}
	if diff := cmp.Diff(expectedSlackIDs, slackIDs); diff != "" {
		t.Errorf("mismatch in Slack IDs (-want +got):\n%s", diff)
	}
}

func TestGetSlackIdsOfCommitters_SlackRateLimit(t *testing.T) {
	config := Config{
		AccessToken:          "mock-access-token",
		CommitterListGitPath: "/mock/repo/path",
	}

	plugin := Plugin{Config: config}
	mockGetAuthorsList := func(gitDir string) ([]string, error) {
		return []string{"user1@example.com", "user2@example.com"}, nil
	}

	mockGetSlackUserIDByEmail := func(accessToken string, emailList string) ([]string, error) {
		return nil, fmt.Errorf("rate_limit_exceeded: too many requests")
	}

	slackIDs, err := GetSlackIdsOfCommitters(&plugin, mockGetAuthorsList, mockGetSlackUserIDByEmail)

	assert.ErrorContains(t, err, "rate_limit_exceeded")
	expectedSlackIDs := []string{}
	if diff := cmp.Diff(expectedSlackIDs, slackIDs); diff != "" {
		t.Errorf("mismatch in Slack IDs (-want +got):\n%s", diff)
	}
}
