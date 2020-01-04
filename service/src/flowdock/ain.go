package flowdock

import (
	"bytes"
	"html"
	"time"
)

const (
	ainMessageTemplate = `  {
    "flow_token": "{{.FlowToken}}",
    "event": "activity",
    "author": {
      "name": "AID bot"
    },
    "title": "New AI Notification",
    "external_thread_id": "AID {{.CurrTime}} {{.AppID}} {{.MessageType}}",
    "thread": {
      "title": "{{.AppID}} {{.MessageType}}",
      "body": "{{.RenderedMessage}}",
      "fields": [
        { "label": "appID", "value": "{{.AppID}}" },
        { "label": "type", "value": "{{.MessageType}}" },
        { "label": "message", "value": "{{.EscapedMessage}}" }
      ],
      "status": {
        "color": "red",
        "value": "Notification"
      }
    }
  }`
)

func (c *Client) buildAiNotificationMessage(appID int32, messageType string, renderedMsg string) (*bytes.Buffer, error) {
	variables := struct {
		FlowToken       string
		CurrTime        time.Time
		AppID           int32
		MessageType     string
		EscapedMessage  string
		RenderedMessage string
	}{
		c.FlowdockToken,
		time.Now().UTC(),
		appID,
		messageType,
		html.EscapeString(renderedMsg),
		c.HTMLReplacer.Replace(renderedMsg),
	}

	return buildMessage(ainMessageTemplate, variables)
}

// SendAiNotificationMessage sends an AI Notification to AID flowdock inbox
func (c *Client) SendAiNotificationMessage(appID int32, messageType string, renderedMsg string) error {
	if c.FlowdockToken == "" {
		return nil
	}
	msg, err := c.buildAiNotificationMessage(appID, messageType, renderedMsg)
	if err != nil {
		return err
	}
	return c.sendMessage(msg)
}
