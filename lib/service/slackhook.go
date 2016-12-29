package service

import (
	"fmt"
	"strings"

	"../common"
	"../helper"
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

// WebhookData used to send a webhook request
// TODO Add formatted HTML and Text Message along with JSON data as Message
type WebhookData struct {
	Referral string
	Message  map[string][]string
	FormName string
}

func (pr *ProcessedRequest) SendToWebhook(webhookURL string) error {
	data := &WebhookData{
		Referral: pr.IncomingRequest.Referral.String(),
		FormName: pr.SingleFormConfig.Name,
		Message:  pr.IncomingRequest.Message,
	}
	return helper.PostJSONtoURL(webhookURL, data)
}

func (pr *ProcessedRequest) SendToSlack(slackURL string) error {
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
		Username:    common.Config.SlackUserName,
		Text:        fmt.Sprintf("*Got a New Submission from %s *:", pr.IncomingRequest.Referral.String()),
		Attachments: attachments,
	}

	return helper.PostJSONtoURL(slackURL, message)
}
