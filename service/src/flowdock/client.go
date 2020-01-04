package flowdock

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
	"time"
)

const (
	flowdockMessagesURI = "https://api.flowdock.com/messages"
)

// Client takes care of flowdock integration
type Client struct {
	FlowdockToken string
	HTTPClient    *http.Client
	HTMLReplacer  *strings.Replacer
}

// NewClient returns a new Flowdock client
func NewClient(flowdockToken string) *Client {
	timeout := time.Duration(5 * time.Second)
	httpClient := http.Client{
		Timeout: timeout,
	}
	return &Client{
		FlowdockToken: flowdockToken,
		HTTPClient:    &httpClient,
		HTMLReplacer:  strings.NewReplacer("\"", "'", "\\n", "<br>"),
	}
}

func buildMessage(mTemplate string, variables interface{}) (*bytes.Buffer, error) {
	tmpl, err := template.New("message").Parse(mTemplate)
	if err != nil {
		return nil, err
	}

	message := new(bytes.Buffer)
	err = tmpl.Execute(message, variables)
	return message, err
}

func (c *Client) sendMessage(message *bytes.Buffer) error {
	request, err := http.NewRequest("POST", flowdockMessagesURI, message)
	if err != nil {
		return err
	}
	request.Header.Set("Content-type", "application/json")

	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	respBodyStr := string(respBody)
	if respBodyStr != "{}" {
		return errors.New(respBodyStr)
	}

	return nil
}
