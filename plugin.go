package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/object"
	"log"
	"os"
	"path/filepath"
	"strings"
	textTemplate "text/template"
	"time"

	"github.com/drone/drone-template-lib/template"
	"github.com/slack-go/slack"

	"github.com/go-git/go-git/v5"
)

type (
	Repo struct {
		Owner string
		Name  string
	}

	BlockSet struct {
		Blocks []json.RawMessage `json:"blocks"`
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
		Webhook        string
		Channel        string
		Recipient      string
		Username       string
		Template       string
		Fallback       string
		ImageURL       string
		IconURL        string
		IconEmoji      string
		Color          string
		LinkNames      bool
		CustomBlock    string
		AccessToken    string
		Mentions       string
		CustomTemplate string
		Message        string
		// File Upload attributes
		FilePath       string
		FileName       string
		InitialComment string
		Title          string
		FailOnError    bool
		// Get Slack ID of the user by email
		SlackIdOf string
		// Git path to get list of committer emails
		CommitterListGitPath string
		CommitterSlackId     bool
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
	var blocks []slack.Block
	var channel string
	var text string
	var fallbackText string

	if p.Config.FilePath != "" {
		return p.UploadFile()
	}

	if p.Config.SlackIdOf != "" {
		return GetSlackIdFromEmail(&p)
	}

	if p.Config.CommitterSlackId && p.Config.Channel == "" {
		_, err := GetSlackIdsOfCommitters(&p, GetChangesetAuthorsList, getSlackUserIDByEmail)
		return err
	}

	// Determine the channel
	if p.Config.Recipient != "" {
		channel = prepend("@", p.Config.Recipient)
	} else if p.Config.Channel != "" {
		channel = prepend("#", p.Config.Channel)
	}

	// Determine the message and fallback
	if p.Config.Template != "" {
		var err error
		text, err = templateMessage(p.Config.Template, p)
		if err != nil {
			return err
		}
	} else if p.Config.Message != "" {
		text = p.Config.Message
	} else {
		text = message(p.Repo, p.Build)
	}

	// Add mentions to the message
	if p.Config.Mentions != "" {
		var mentionUserIDs = strings.Split(p.Config.Mentions, ",")
		mentions := make([]string, len(mentionUserIDs))
		for i, id := range mentionUserIDs {
			// Check if the id starts with "@" and format it accordingly
			if strings.HasPrefix(id, "@") {
				mentions[i] = fmt.Sprintf("<%s>:", id)
			} else {
				mentions[i] = fmt.Sprintf("<@%s>:", id)
			}
		}
		mentionText := strings.Join(mentions, " ")
		text = fmt.Sprintf("%s %s", mentionText, text)
	}
	if p.Config.CustomTemplate != "" {
		// Read JSON from file
		var filePath string

		switch p.Config.CustomTemplate {
		case "basic_success_1":
			filePath = "templates/basic_success.json"
		case "basic_fail_1":
			filePath = "templates/basic_fail.json"
		case "success_tagged_deploy_1":
			filePath = "templates/success_tag_deploy.json"
		case "basic_on_hold_1":
			filePath = "templates/basic_on_hold.json"
		default:
			return fmt.Errorf("invalid template name: %s", p.Config.CustomTemplate)
		}

		// Read JSON from file
		file, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read template file: %w", err)
		}

		// Fill in the missing values in the template
		tmpl, err := textTemplate.New("template").Parse(string(file))
		if err != nil {
			return fmt.Errorf("failed to parse template: %w", err)
		}

		var filledTemplate bytes.Buffer
		err = tmpl.Execute(&filledTemplate, p)
		if err != nil {
			return fmt.Errorf("failed to fill in template values: %w", err)
		}

		// Parse the filled template JSON into a BlockSet
		var blockSet BlockSet
		err = json.Unmarshal(filledTemplate.Bytes(), &blockSet)
		if err != nil {
			return fmt.Errorf("failed to parse filled template JSON: %w", err)
		}

		// Parse each block into a slack.Block and append to blocks
		for _, rawBlock := range blockSet.Blocks {
			var blockType struct {
				Type string `json:"type"`
			}
			err := json.Unmarshal(rawBlock, &blockType)
			if err != nil {
				return fmt.Errorf("failed to parse block type JSON: %w", err)
			}

			var block slack.Block
			switch blockType.Type {
			case "section":
				block = new(slack.SectionBlock)
			case "divider":
				block = new(slack.DividerBlock)
			case "header":
				block = new(slack.HeaderBlock)
			case "actions":
				block = new(slack.ActionBlock)
			default:
				return fmt.Errorf("unknown block type: %s", blockType.Type)
			}

			err = json.Unmarshal(rawBlock, block)
			if err != nil {
				return fmt.Errorf("failed to parse block JSON: %w", err)
			}

			blocks = append(blocks, block)
		}
		text = ""
	}

	if p.Config.Fallback != "" {
		var err error
		fallbackText, err = templateMessage(p.Config.Fallback, p)
		if err != nil {
			return err
		}
	} else {
		fallbackText = fallback(p.Repo, p.Build)
	}

	// Determine the color
	colorText := p.Config.Color
	if colorText == "" {
		colorText = color(p.Build)
	}

	// Parse custom blocks if they exist
	if p.Config.CustomBlock != "" {
		var blockSet BlockSet
		err := json.Unmarshal([]byte(p.Config.CustomBlock), &blockSet)
		if err != nil {
			return fmt.Errorf("could not unmarshal custom block: %w", err)
		}
		for _, rawMessage := range blockSet.Blocks {
			block := new(slack.SectionBlock)
			err := json.Unmarshal(rawMessage, block)
			if err != nil {
				return fmt.Errorf("could not unmarshal individual block: %w", err)
			}
			blocks = append(blocks, block)
		}
	}

	// If access token is provided, use it
	if p.Config.AccessToken != "" {
		slackApi := slack.New(p.Config.AccessToken)
		_, err := slackApi.AuthTest()
		if err != nil {
			return fmt.Errorf("failed to authenticate using access token: %w", err)
		}

		options := []slack.MsgOption{}
		if len(blocks) > 0 {
			options = append(options, slack.MsgOptionBlocks(blocks...))
		} else {
			options = append(options, slack.MsgOptionText(text, false))
		}

		_, _, err = slackApi.PostMessage(channel, options...)
		if err != nil {
			return fmt.Errorf("failed to post message using access token: %w", err)
		}

		if p.Config.CommitterSlackId {
			err := p.sendDirectMessageToCommitters(options)
			if err != nil {
				return fmt.Errorf("failed to send direct message to committers: %w", err)
			}
		}

		return nil
	}

	// Build the attachment
	attachment := slack.Attachment{
		Color:      colorText,
		ImageURL:   p.Config.ImageURL,
		MarkdownIn: []string{"text", "fallback"},
		Text:       text,
		Fallback:   fallbackText,
	}

	// Build the payload
	payload := slack.WebhookMessage{
		Username:    p.Config.Username,
		Attachments: []slack.Attachment{attachment},
		IconURL:     p.Config.IconURL,
		IconEmoji:   p.Config.IconEmoji,
		Channel:     channel,
	}

	// Add custom blocks to the payload if they exist
	if len(blocks) > 0 {
		payload.Blocks = &slack.Blocks{
			BlockSet: blocks,
		}
	}

	// Post the message with the webhook
	return slack.PostWebhook(p.Config.Webhook, &payload)
}

func (p Plugin) UploadFile() error {

	p.Config.FilePath = strings.TrimSpace(p.Config.FilePath)

	api := slack.New(p.Config.AccessToken)
	fileSize, err := GetFileSize(p.Config.FilePath)
	if err != nil {
		log.Printf("Error getting file size: %s\n", err.Error())
		return err
	}

	if p.Config.FileName == "" {
		fileName := filepath.Base(p.Config.FilePath)
		p.Config.FileName = fileName
	}

	params := slack.UploadFileV2Parameters{
		File:           p.Config.FilePath,
		Channel:        p.Config.Channel,
		Filename:       p.Config.FileName,
		Title:          p.Config.Title,
		InitialComment: p.Config.InitialComment,
		FileSize:       fileSize,
	}

	slackSummary, err := api.UploadFileV2(params)

	if !p.Config.FailOnError && slackSummary == nil {
		if err != nil {
			log.Println("Bad Api ret val, upload file failed but passing build as PLUGIN_FAIL_ON_ERROR is false")
		}
		return nil
	} else if p.Config.FailOnError && slackSummary == nil {
		log.Println("Bad ret val,  Failed to upload file, failing build")
		_ = p.WriteFileUploadResult("", "", err)
		return fmt.Errorf("Bad ret val, Failed to upload file %s ", p.Config.FilePath)
	}

	if !p.Config.FailOnError && err != nil {
		if err != nil {
			log.Println("Unable to upload file but passing build PLUGIN_FAIL_ON_ERROR is false")
		}
		return nil
	} else if p.Config.FailOnError && err != nil {
		log.Println("Upload API Failed to upload file, failing build")
		_ = p.WriteFileUploadResult("", "", err)
		return fmt.Errorf("Failed to upload file %s ", p.Config.FilePath)
	}

	err = p.WriteFileUploadResult(slackSummary.ID, slackSummary.Title, err)
	if !p.Config.FailOnError {
		if err != nil {
			log.Println("Unable to Write output env var results for file upload " +
				"but passing build PLUGIN_FAIL_ON_ERROR is false")
		}
		return nil
	}

	return nil
}

func (p Plugin) WriteFileUploadResult(slackFileId, title string, err error) error {

	type EnvKvPair struct {
		Key   string
		Value string
	}

	resultStr := "Failed: Slack file upload failed"
	if err == nil {
		resultStr = "Success: Slack file upload successful"
	}

	var kvPairs = []EnvKvPair{
		{Key: "UPLOAD_OK_STATUS", Value: resultStr},
		{Key: "UPLOAD_FILE_PATH", Value: p.Config.FilePath},
	}

	var retErr error = nil

	for _, kvPair := range kvPairs {
		err := WriteEnvToOutputFile(kvPair.Key, kvPair.Value)
		if err != nil {
			retErr = err
		}
	}

	return retErr
}

func WriteEnvToOutputFile(key, value string) error {
	outputFile, err := os.OpenFile(os.Getenv("DRONE_OUTPUT"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer outputFile.Close()
	_, err = fmt.Fprintf(outputFile, "%s=%s\n", key, value)
	if err != nil {
		return fmt.Errorf("failed to write to env: %w", err)
	}
	return nil
}

func GetFileSize(filePath string) (int, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to get file info: %w", err)
	}
	if fileInfo.IsDir() {
		return 0, fmt.Errorf("path %s is a directory, not a file", filePath)
	}
	return int(fileInfo.Size()), nil
}

func templateMessage(t string, plugin Plugin) (string, error) {
	c, err := contents(t)
	if err != nil {
		return "", fmt.Errorf("could not read template: %w", err)
	}

	return template.RenderTrim(c, plugin)
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

func GetSlackIdsOfCommitters(p *Plugin,
	getAuthorsListFunc func(string) ([]string, error),
	getSlackUserIDByEmailFunc func(string, string) ([]string, error)) ([]string, error) {

	if p.Config.CommitterListGitPath == "" {
		p.Config.CommitterListGitPath = os.Getenv("DRONE_WORKSPACE")
	}

	emails, err := getAuthorsListFunc(p.Config.CommitterListGitPath)
	if err != nil {
		log.Println("Failed to get git emails: ", err)
		return []string{}, fmt.Errorf("failed to get git emails: %w", err)
	}

	slackUserIdList, err := getSlackUserIDByEmailFunc(p.Config.AccessToken, strings.Join(emails, ","))
	if err != nil {
		log.Println("Failed to get Slack ID by email: ", err)
		return []string{}, fmt.Errorf("failed to get Slack ID by email: %w", err)
	}

	jsonStr := strings.Join(slackUserIdList, ",")
	err = WriteEnvToOutputFile("COMMITTERS_SLACK_IDS", jsonStr)
	if err != nil {
		log.Println("Failed to write git emails to output file: ", err)
		return []string{}, fmt.Errorf("failed to write git emails to output file: %w", err)
	}

	return slackUserIdList, nil
}

func GetSlackIdFromEmail(p *Plugin) error {
	slackIdList, err := getSlackUserIDByEmail(p.Config.AccessToken, p.Config.SlackIdOf)
	if err != nil {
		log.Println("Failed to get Slack ID by email: ", err)
		return fmt.Errorf("failed to get Slack ID by email: %w", err)
	}

	slackIdsCsvStr := strings.Join(slackIdList, ",")
	err = WriteEnvToOutputFile("SLACK_ID_FROM_EMAIL", slackIdsCsvStr)
	if err != nil {
		return fmt.Errorf("failed to write Slack ID to output file: %w", err)
	}
	return nil
}

func getSlackUserIDByEmail(accessToken, emailListStr string) ([]string, error) {

	emailArray := []string{}
	for _, email := range strings.Split(emailListStr, ",") {
		trimmedEmail := strings.TrimSpace(email)
		if trimmedEmail != "" {
			emailArray = append(emailArray, trimmedEmail)
		}
	}
	slackIdsList := []string{}

	var failedEmails []string
	for _, email := range emailArray {
		api := slack.New(accessToken)
		if api == nil {
			log.Println("Failed to create Slack client")
			return emailArray, fmt.Errorf("failed to create Slack client")
		}

		user, err := api.GetUserByEmail(email)
		if err != nil {
			log.Printf("Failed to fetch Slack ID for email %s: %v", email, err)
			failedEmails = append(failedEmails, email)
			continue
		}
		slackIdsList = append(slackIdsList, user.ID)

		// Add a short delay to avoid rate limits
		time.Sleep(500 * time.Millisecond)
	}
	if len(failedEmails) > 0 {
		log.Printf("Failed to fetch Slack IDs for the following emails: %v", failedEmails)
	}

	return slackIdsList, nil
}

func (p Plugin) sendDirectMessageToCommitters(options []slack.MsgOption) error {
	slackUserIdList, err := GetSlackIdsOfCommitters(&p, GetChangesetAuthorsList, getSlackUserIDByEmail)
	if err != nil {
		log.Println("Failed to get Slack ID by email: ", err)
		return fmt.Errorf("failed to get Slack ID by email: %w", err)
	}
	for _, slackUserId := range slackUserIdList {
		err = sendDirectMessage(p.Config.AccessToken, slackUserId, options)
		if err != nil {
			log.Println("Failed to send direct message: ", err)
			continue
		}
		log.Println("Message sent successfully for ", slackUserId)
	}
	return nil
}

func sendDirectMessage(botToken, userID string, options []slack.MsgOption) error {

	client := slack.New(botToken)

	channel, _, _, err := client.OpenConversation(&slack.OpenConversationParameters{
		ReturnIM: true,
		Users:    []string{userID},
	})
	if err != nil {
		log.Println("Failed to open conversation: ", err)
	}

	_, _, err = client.PostMessage(channel.ID, options...)
	if err != nil {
		log.Printf("Failed to send direct slack message: %v", err)
	}

	return nil
}

func GetChangesetAuthorsList(gitDir string) ([]string, error) {
	if gitDir == "" {
		log.Println("gitDir is empty")
		return nil, fmt.Errorf("gitDir cannot be empty")
	}

	absGitDir, err := filepath.Abs(gitDir)
	if err != nil {
		log.Println("Failed to get absolute path of gitDir: ", gitDir)
		return nil, fmt.Errorf("failed to get absolute path of gitDir: %w", err)
	}

	repo, err := git.PlainOpen(absGitDir)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository: %w", err)
	}

	headRef, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	headCommit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD commit: %w", err)
	}

	parentCommitIter := headCommit.Parents()
	oldCommit, err := parentCommitIter.Next()
	if err != nil {
		return nil, fmt.Errorf("failed to get parent commit of HEAD: %w", err)
	}

	emailSet := make(map[string]struct{})

	commitIter, err := repo.Log(&git.LogOptions{
		From:  headCommit.Hash,
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get commit log: %w", err)
	}

	err = commitIter.ForEach(func(commit *object.Commit) error {
		// Stop if we reach the parent commit
		if commit.Hash == oldCommit.Hash {
			return nil
		}

		email := strings.TrimSpace(commit.Author.Email)
		if email != "" {
			emailSet[email] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error during commit log iteration: %w", err)
	}

	var uniqueEmails []string
	for email := range emailSet {
		uniqueEmails = append(uniqueEmails, email)
	}
	return uniqueEmails, nil
}
