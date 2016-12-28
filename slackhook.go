package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// SlackMessage ..
type SlackMessage struct {
	Username    string            `json:"username"`
	Text        string            `json:"text"`
	Attachments []SlackAttachment `json:"attachments"`
}

// SlackAttachment ..
type SlackAttachment struct {
	Fallback       string                 `json:"fallback"`
	Title          string                 `json:"title"`
	Color          string                 `json:"color,omitempty"`
	PreText        string                 `json:"pretext"`
	AuthorName     string                 `json:"author_name"`
	AuthorLink     string                 `json:"author_link"`
	Fields         []SlackAttachmentField `json:"fields"`
	MarkdownFormat []string               `json:"mrkdwn_in"`
	Text           string                 `json:"text"`
	ThumbnailURL   string                 `json:"thumb_url,omitempty"`
}

// SlackAttachmentField ..
type SlackAttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// SlackTypeLink ..
type SlackTypeLink struct {
	Text string
	Href string
}

// <http://www.amazon.com|Amazon>
func (s *SlackTypeLink) String() string {
	return fmt.Sprintf("<%s|%s>", s.Href, s.Text)
}

// move to helper
func postJSONtoURL(url string, data interface{}) error {
	pr, pw := io.Pipe()
	go func() {
		// close the writer, so the reader knows there's no more data
		defer pw.Close()

		// write json data to the PipeReader through the PipeWriter
		if err := json.NewEncoder(pw).Encode(data); err != nil {
			log.Print(err)
		}
	}()

	if _, err := http.Post(url, "application/json", pr); err != nil {
		return err
	}
	return nil
}

// WebhookData used to send a webhook request
// TODO Add formatted HTML and Text Message along with JSON data as Message
type WebhookData struct {
	Referral string
	Message  map[string][]string
	FormName string
}

func (pr *ProcessedRequest) sendToWebhook(webhookURL string) error {
	data := &WebhookData{
		Referral: pr.IncomingRequest.Referral.String(),
		FormName: pr.SingleFormConfig.Name,
		Message:  pr.IncomingRequest.Message,
	}
	return postJSONtoURL(webhookURL, data)
}

func (pr *ProcessedRequest) sendToSlack(slackURL string) error {

	attachments := make([]SlackAttachment, 1)

	msg := pr.IncomingRequest.Message
	var d string
	fields := make([]SlackAttachmentField, 0, len(msg))
	for key, data := range msg {
		if len(data) == 0 {
			continue
		} else if len(data) > 1 {
			d = strings.Join(data, "; ")
		} else {
			d = data[0]
		}
		field := SlackAttachmentField{
			Title: key,
			Value: d,
			Short: false,
		}
		fields = append(fields, field)
	}

	attachment := SlackAttachment{
		Title:          pr.SingleFormConfig.Name,
		Fields:         fields,
		MarkdownFormat: []string{"fields"},
	}
	attachments = append(attachments, attachment)

	message := &SlackMessage{
		Username:    config.SlackUserName,
		Text:        fmt.Sprintf("*Got a New Submission from %s *:", pr.IncomingRequest.Referral.String()),
		Attachments: attachments,
	}

	return postJSONtoURL(slackURL, message)
}
