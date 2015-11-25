package slack

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Client interface {
	SendMessage(*Message) error
}

type client struct {
	url string
}

func NewClient(url string) Client {
	return &client{url}
}

func (c *client) SendMessage(msg *Message) error {

	body, _ := json.Marshal(msg)
	buf := bytes.NewReader(body)

	resp, err := http.Post(c.url, "application/json", buf)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		t, _ := ioutil.ReadAll(resp.Body)
		return &Error{resp.StatusCode, string(t)}
	}

	return nil
}
