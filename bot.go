package gozulipbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
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

func (b Bot) GetStreamList() (*http.Response, error) {
	req, err := b.constructRequest("GET", "streams", "")
	if err != nil {
		return nil, err
	}

	c := http.Client{}
	return c.Do(req)
}

type streamJson struct {
	Msg     string `json:msg`
	Streams []struct {
		StreamID    int    `json:stream_id`
		InviteOnly  bool   `json:invite_only`
		Description string `json:description`
		Name        string `json:name`
	} `json:streams`
	Result string `json:result`
}

func (b Bot) GetStreamNameList() ([]string, error) {
	resp, err := b.GetStreamList()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var getStreamJson streamJson

	err = json.Unmarshal(body, &getStreamJson)
	if err != nil {
		return nil, err
	}

	var outStreams []string

	for _, stream := range getStreamJson.Streams {
		outStreams = append(outStreams, stream.Name)
	}

	return outStreams, nil
}

func (b Bot) SubscribeToStreams(streams []string) (*http.Response, error) {
	// TODO subscribe to streams the bot has if no streams given
	var toSubStreams []map[string]string
	for _, name := range streams {
		toSubStreams = append(toSubStreams, map[string]string{"name": name})
	}

	bodyBts, err := json.Marshal(toSubStreams)
	if err != nil {
		return nil, err
	}

	body := "subscriptions=" + string(bodyBts)

	req, err := b.constructRequest("POST", "users/me/subscriptions", body)
	if err != nil {
		return nil, err
	}

	c := http.Client{}
	return c.Do(req)
}

func (b Bot) RegisterEvents() (*http.Response, error) {
	req, err := b.constructRequest("POST", "register", `event_types=["message"]`)
	if err != nil {
		return nil, err
	}

	c := http.Client{}
	return c.Do(req)
}

func (b Bot) GetEventsFromQueue(queueID string,
	lastMessageID int) (*http.Response, error) {
	values := url.Values{}
	values.Set("queue_id", queueID)
	values.Set("last_event_id", strconv.Itoa(lastMessageID))

	url := "events?" + values.Encode()

	req, err := b.constructRequest("GET", url, "")
	if err != nil {
		return nil, err
	}

	c := http.Client{}
	return c.Do(req)
}

func (b Bot) RespondToMessage(e EventMessage, response string) (*http.Response,
	error) {
	if response == "" {
		return nil, fmt.Errorf("Message response cannot be blank")
	}
	if e.DisplayRecipient.Topic != "" {
		return b.SendStreamMessage(e.DisplayRecipient.Topic, e.Subject, response)
	}
	// TODO handle multiple users in a private message
	if e.Subject == "" {
		return b.SendPrivateMessage(e.SenderEmail, response)
	}
	return nil, fmt.Errorf("EventMessage is not understood: %v\n", e)
}

func (b Bot) RespondToMessagePrivately(e EventMessage,
	response string) (*http.Response, error) {
	if response == "" {
		return nil, fmt.Errorf("Message response cannot be blank")
	}
	return b.SendPrivateMessage(e.SenderEmail, response)
}

func (b Bot) constructRequest(method, endpoint, body string) (*http.Request,
	error) {
	url := fmt.Sprintf("https://api.zulip.com/v1/%s", endpoint)
	req, err := http.NewRequest(method, url, strings.NewReader(body))
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

	return b.constructRequest("POST", "messages", values.Encode())
}
