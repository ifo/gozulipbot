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
	Email   string
	APIKey  string
	Queues  []*Queue
	Streams []string
	client  Doer
}

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type Queue struct {
	ID           string
	LastEventID  int
	MaxMessageID int
}

// Init adds an http client to an existing bot struct.
func (b *Bot) Init() *Bot {
	b.client = &http.Client{}
	return b
}

// GetStreamList gets the raw http response when requesting all public streams.
func (b *Bot) GetStreamList() (*http.Response, error) {
	req, err := b.constructRequest("GET", "streams", "")
	if err != nil {
		return nil, err
	}

	return b.client.Do(req)
}

type StreamJSON struct {
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

	var sj StreamJSON
	err = json.Unmarshal(body, &sj)
	if err != nil {
		return nil, err
	}

	var streams []string
	for _, s := range sj.Streams {
		streams = append(streams, s.Name)
	}

	return streams, nil
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

func (b *Bot) ListSubscriptions() (*http.Response, error) {
	req, err := b.constructRequest("GET", "users/me/subscriptions", "")
	if err != nil {
		return nil, err
	}

	return b.client.Do(req)
}

type EventType string

const (
	Messages      EventType = "messages"
	Subscriptions EventType = "subscriptions"
	RealmUser     EventType = "realm_user"
	Pointer       EventType = "pointer"
)

type Narrow string

const (
	NarrowPrivate Narrow = `[["is", "private"]]`
	NarrowAt      Narrow = `[["is", "mentioned"]]`
)

// RegisterEvents tells Zulip to include message events in the bots events queue.
// Passing nil as the slice of EventType will default to receiving Messages
func (b *Bot) RegisterEvents(es []EventType, n Narrow) (*http.Response, error) {
	// default to Messages if no EventTypes given
	query := `event_types=["message"]`

	if len(es) != 0 {
		query = `event_types=["`
		for i, s := range es {
			query += fmt.Sprintf("%s", s)
			if i != len(es)-1 {
				query += `", "`
			}
		}
		query += `"]`
	}

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
	return b.RegisterEvents(nil, "")
}

func (b *Bot) RegisterAt() (*http.Response, error) {
	return b.RegisterEvents(nil, NarrowAt)
}

func (b *Bot) RegisterPrivate() (*http.Response, error) {
	return b.RegisterEvents(nil, NarrowPrivate)
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

// constructRequest makes a zulip request and ensures the proper headers are set.
func (b *Bot) constructRequest(method, endpoint, body string) (*http.Request, error) {
	url := "https://api.zulip.com/v1/" + endpoint
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(b.Email, b.APIKey)

	return req, nil
}
