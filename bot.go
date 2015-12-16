package gozulipbot

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
	req, err := constructRequest(b, "stream", stream, topic, content)
	if err != nil {
		return nil, err
	}

	c := http.Client{}
	return c.Do(req)
}

func (b Bot) SendPrivateMessage(email, content string) (*http.Response, error) {
	// TODO ensure "user" (a.k.a. email) exists, content is non-empty
	req, err := constructRequest(b, "private", email, "", content)
	if err != nil {
		return nil, err
	}

	c := http.Client{}
	return c.Do(req)
}

func constructRequest(bot Bot, mtype, to, subject,
	content string) (*http.Request, error) {
	v := url.Values{}
	v.Set("type", mtype)
	v.Set("to", to)
	v.Set("content", content)
	if mtype == "stream" {
		v.Set("subject", subject)
	}

	req, err := http.NewRequest("POST", "https://api.zulip.com/v1/messages",
		strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(bot.EmailAddress, bot.ApiKey)

	return req, nil
}
