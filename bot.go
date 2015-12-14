package main

import (
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
	v := url.Values{}
	v.Set("type", "stream")
	v.Set("to", stream)
	v.Set("subject", topic)
	v.Set("content", content)

	req, err := http.NewRequest("POST", "https://api.zulip.com/v1/messages",
		strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(b.EmailAddress, b.ApiKey)

	c := http.Client{}
	return c.Do(req)
}

func (b Bot) SendPrivateMessage(email, content string) (*http.Response, error) {
	// TODO ensure "user" (a.k.a. email) exists, content is non-empty
	v := url.Values{}
	v.Set("type", "private")
	v.Set("to", email)
	v.Set("content", content)

	req, err := http.NewRequest("POST", "https://api.zulip.com/v1/messages",
		strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(b.EmailAddress, b.ApiKey)

	c := http.Client{}
	return c.Do(req)
}
