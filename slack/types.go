package slack

import "fmt"

type Message struct {
	Text        string        `json:"text"`
	Username    string        `json:"username"`
	IconUrl     string        `json:"icon_url"`
	IconEmoji   string        `json:"icon_emoji"`
	Channel     string        `json:"channel"`
	UnfurlLinks bool          `json:"unfurl_links"`
	Attachments []*Attachment `json:"attachments"`
}

func (m *Message) NewAttachment() *Attachment {
	a := &Attachment{}
	m.AddAttachment(a)
	return a
}

func (m *Message) AddAttachment(a *Attachment) {
	m.Attachments = append(m.Attachments, a)
}

type Attachment struct {
	Fallback string   `json:"fallback"`
	Text     string   `json:"text"`
	Pretext  string   `json:"pretext"`
	Color    string   `json:"color"`
	Fields   []*Field `json:"fields"`
	MrkdwnIn []string `json:"mrkdwn_in"`
}

func (a *Attachment) NewField() *Field {
	f := &Field{}
	a.AddField(f)
	return f
}

func (a *Attachment) AddField(f *Field) {
	a.Fields = append(a.Fields, f)
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type Error struct {
	Code int
	Body string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Slack Error: %d %s", e.Code, e.Body)
}
