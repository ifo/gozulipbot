package gozulipbot

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type Queue struct {
	ID           string `json:"queue_id"`
	LastEventID  int    `json:"last_event_id"`
	MaxMessageID int    `json:"max_message_id"`
	Bot          *Bot   `json:"-"`
}

// GetEvents is a blocking call that waits for and parses a list of EventMessages.
// There will usually only be one EventMessage returned.
// When a heartbeat is returned, GetEvents will return a HeartbeatError
func (q *Queue) GetEvents() ([]EventMessage, error) {
	resp, err := q.RawGetEvents()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	msgs, err := ParseEventMessages(body)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

// RawGetEvents is a blocking call that receives a response containing a list
// of events (a.k.a. received messages) since the last message id in the queue.
func (q *Queue) RawGetEvents() (*http.Response, error) {
	values := url.Values{}
	values.Set("queue_id", q.ID)
	values.Set("last_event_id", strconv.Itoa(q.LastEventID))

	url := "events?" + values.Encode()

	req, err := q.Bot.constructRequest("GET", url, "")
	if err != nil {
		return nil, err
	}

	return q.Bot.Client.Do(req)
}
