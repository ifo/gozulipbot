package gozulipbot

import (
	"encoding/json"
	"errors"
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
	client       Doer
}

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// A Message is all of the necessary metadata to post on Zulip.
// It can be either a public message, where Topic is set, or a private message,
// where there is at least one element in Emails.
//
// If the length of Emails is not 0, functions will always assume it is a private message.
type Message struct {
	Stream  string
	Topic   string
	Emails  []string
	Content string
}

// MakeBot creates a bot object and gives it an http client.
func MakeBot(email, apikey string, streams []string) Bot {
	return Bot{
		EmailAddress: email,
		ApiKey:       apikey,
		Streams:      streams,
		client:       http.DefaultClient,
	}
}

// Message posts a message to Zulip. If any emails have been set on the message,
// the message will be re-routed to the PrivateMessage function.
func (b *Bot) Message(m Message) (*http.Response, error) {
	if m.Content == "" {
		return nil, errors.New("content cannot be empty")
	}

	// if any emails are set, this is a private message
	if len(m.Emails) != 0 {
		return b.PrivateMessage(m)
	}

	// otherwise it's a stream message
	if m.Stream == "" {
		return nil, errors.New("stream cannot be empty")
	}
	if m.Topic == "" {
		return nil, errors.New("topic cannot be empty")
	}
	req, err := b.constructMessageRequest(m)
	if err != nil {
		return nil, err
	}
	return b.client.Do(req)
}

// PrivateMessage sends a message to the first user in the message email slice.
func (b *Bot) PrivateMessage(m Message) (*http.Response, error) {
	if len(m.Emails) == 0 {
		return nil, errors.New("there must be at least one recipient")
	}
	req, err := b.constructMessageRequest(m)
	if err != nil {
		return nil, err
	}

	return b.client.Do(req)
}

// GetStreamList gets the raw http response when requesting all public streams.
func (b *Bot) GetStreamList() (*http.Response, error) {
	req, err := b.constructRequest("GET", "streams", "")
	if err != nil {
		return nil, err
	}

	return b.client.Do(req)
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

// GetStreams returns a list of all public streams
func (b *Bot) GetStreams() ([]string, error) {
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

// Subscribe will set the bot to receive messages from the given streams.
// If no streams are given, it will subscribe the bot to the streams in the bot struct.
func (b *Bot) Subscribe(streams []string) (*http.Response, error) {
	if streams == nil {
		streams = b.Streams
	}

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

	return b.client.Do(req)
}

type Narrow string

const (
	NarrowPrivate Narrow = `[["is", "private"]]`
	NarrowAt      Narrow = `[["is", "mentioned"]]`
)

// RegisterEvents tells Zulip to include message events in the bots events queue.
// It is necessary to call only once ever, to be able to receive messages.
// Calling it multiple times will have no negative effect.
func (b *Bot) RegisterEvents(evtTypes []string, n Narrow) (*http.Response, error) {
	query := `event_types=["`
	for i, s := range evtTypes {
		query += s
		if i != len(evtTypes)-1 {
			query += `", "`
		}
	}
	query += `"]`

	if n != "" {
		query += fmt.Sprintf("&narrow=%s", n)
	}

	req, err := b.constructRequest("POST", "register", query)
	if err != nil {
		return nil, err
	}

	return b.client.Do(req)
}

func (b *Bot) RegisterAll() (*http.Response, error) {
	return b.RegisterEvents([]string{"message"}, "")
}

func (b *Bot) RegisterAt() (*http.Response, error) {
	return b.RegisterEvents([]string{"message"}, NarrowAt)
}

func (b *Bot) RegisterPrivate() (*http.Response, error) {
	return b.RegisterEvents([]string{"message"}, NarrowPrivate)
}

// GetEventsFromQueue receives a list of events (a.k.a. received messages) since
// the last message given.
// Messages received in this queue will be EventMessages.
func (b *Bot) GetEventsFromQueue(queueID string, lastMessageID int) (*http.Response, error) {
	values := url.Values{}
	values.Set("queue_id", queueID)
	values.Set("last_event_id", strconv.Itoa(lastMessageID))

	url := "events?" + values.Encode()

	req, err := b.constructRequest("GET", url, "")
	if err != nil {
		return nil, err
	}

	return b.client.Do(req)
}

// Respond sends a given message as a response to whatever context from which
// an EventMessage was received.
func (b *Bot) Respond(e EventMessage, response string) (*http.Response, error) {
	if response == "" {
		return nil, errors.New("Message response cannot be blank")
	}
	m := Message{
		Stream:  e.DisplayRecipient.Topic,
		Topic:   e.Subject,
		Content: response,
	}
	if m.Topic != "" {
		return b.Message(m)
	}
	// private message
	if m.Stream == "" {
		emails, err := b.privateResponseList(e)
		if err != nil {
			return nil, err
		}
		m.Emails = emails
		return b.Message(m)
	}
	return nil, fmt.Errorf("EventMessage is not understood: %v\n", e)
}

// privateResponseList gets the list of other users in a private multiple
// message conversation.
func (b *Bot) privateResponseList(e EventMessage) ([]string, error) {
	var out []string
	for _, u := range e.DisplayRecipient.Users {
		if u.Email != b.EmailAddress {
			out = append(out, u.Email)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("EventMessage had no Users within the DisplayRecipient")
	}
	return out, nil
}

// constructRequest makes a zulip request and ensures the proper headers are set.
func (b *Bot) constructRequest(method, endpoint, body string) (*http.Request, error) {
	url := "https://api.zulip.com/v1/" + endpoint
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(b.EmailAddress, b.ApiKey)

	return req, nil
}

// constructMessageRequest is a helper for simplifying sending a message.
func (b *Bot) constructMessageRequest(m Message) (*http.Request, error) {
	to := m.Stream
	mtype := "stream"

	le := len(m.Emails)
	if le != 0 {
		mtype = "private"
	}
	if le == 1 {
		to = m.Emails[0]
	}
	if le > 1 {
		for i, e := range m.Emails {
			to += e
			if i != le-1 {
				to += ", "
			}
		}
	}

	values := url.Values{}
	values.Set("type", mtype)
	values.Set("to", to)
	values.Set("content", m.Content)
	if mtype == "stream" {
		values.Set("subject", m.Topic)
	}

	return b.constructRequest("POST", "messages", values.Encode())
}
