package main

import (
	"github.com/drone/drone-plugin-go/plugin"
	"os"
	"testing"
)

type MockClient struct {
	Url     string
	Message *Message
}

func (c *MockClient) SetUrl(url string) {
	c.Url = url
}

func (c *MockClient) SendMessage(msg *Message) error {
	c.Message = msg
	return nil
}

func readFile(name string) *os.File {
	f, _ := os.Open(name)
	return f
}

func TestSucess(t *testing.T) {
	cases := []struct {
		input    *os.File
		url      string
		username string
		channel  string
		text     string
	}{
		{readFile("testdata/basic.json"), "https://hooks.slack.com/services/...", "drone", "#dev", "*success* <http://drone.mycompany.com/foo/bar/22|foo/bar#9f2849d5> (master) by johnsmith"},
		{readFile("testdata/template.json"), "https://hooks.slack.com/services/...", "drone", "#dev", "*success* foo/bar#9f2849d5 (master) took 3m30s"},
	}
	for _, c := range cases {
		plugin.Stdin = plugin.NewParamSet(c.input)
		mock := &MockClient{}
		client = mock
		main()
		if mock.Url != c.url {
			t.Errorf("Incorrect url: %v, expected %v", mock.Url, c.url)
		}
		if mock.Message.Username != c.username {
			t.Errorf("Incorrect username: %v, expected %v", mock.Message.Username, c.username)
		}
		if mock.Message.Channel != c.channel {
			t.Errorf("Incorrect channel: %v, expected %v", mock.Message.Channel, c.username)
		}
		att := mock.Message.Attachments[0]
		if att.Text != c.text {
			t.Errorf("Incorrect text: %v, expected %v", att.Text, c.text)
		}
	}
}
