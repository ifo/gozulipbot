package gozulipbot

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Bot struct {
	EmailAddress string
	ApiKey       string
	Streams      []string
}

func MakeBot(email, apikey string, streams []string) Bot {
	return Bot{
		EmailAddress: email,
		ApiKey:       apikey,
		Streams:      streams,
	}
}

func (b Bot) SendStreamMessage(stream, topic, content string) (*http.Response,
	error) {
	// TODO ensure stream exists, content is non-empty
	req, err := b.constructMessageRequest("stream", stream, topic, content)
	if err != nil {
		return nil, err
	}

	c := http.Client{}
	return c.Do(req)
}

func (b Bot) SendPrivateMessage(email, content string) (*http.Response, error) {
	// TODO ensure "user" (a.k.a. email) exists, content is non-empty
	req, err := b.constructMessageRequest("private", email, "", content)
	if err != nil {
		return nil, err
	}

	c := http.Client{}
	return c.Do(req)
}

func (b Bot) constructRequest(endpoint, body string) (*http.Request, error) {
	url := fmt.Sprintf("https://api.zulip.com/v1/%s", endpoint)
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(b.EmailAddress, b.ApiKey)

	return req, nil
}

func (b Bot) constructMessageRequest(mtype, to, subject,
	content string) (*http.Request, error) {
	values := url.Values{}
	values.Set("type", mtype)
	values.Set("to", to)
	values.Set("content", content)
	if mtype == "stream" {
		values.Set("subject", subject)
	}

	return b.constructRequest("messages", values.Encode())
}
