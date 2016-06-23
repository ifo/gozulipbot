package gozulipbot

import (
	"encoding/json"
	"errors"
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

func (q *Queue) EventsChan() (chan EventMessage, func()) {
	end := false
	endFunc := func() {
		end = true
	}

	out := make(chan EventMessage)
	go func() {
		defer close(out)
		for {
			if end {
				return
			}
			ems, err := q.GetEvents()
			// TODO? do something with the error
			if err != nil {
				continue
			}
			for _, em := range ems {
				out <- em
			}
		}
	}()

	return out, endFunc
}

// EventsCallback will repeatedly call provided callback function with
// the output of continual queue.GetEvents calls.
// It returns a function which can be called to end the calls.
//
// Note, it will never return a HeartbeatError.
func (q *Queue) EventsCallback(fn func(EventMessage, error)) func() {
	end := false
	endFunc := func() {
		end = true
	}
	go func() {
		for {
			if end {
				return
			}
			ems, err := q.GetEvents()
			if err == HeartbeatError {
				continue
			}
			for _, em := range ems {
				fn(em, err)
			}
		}
	}()

	return endFunc
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

	msgs, err := q.ParseEventMessages(body)
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

var HeartbeatError = errors.New("EventMessage is a heartbeat")

func (q *Queue) ParseEventMessages(rawEventResponse []byte) ([]EventMessage, error) {
	rawResponse := map[string]json.RawMessage{}
	err := json.Unmarshal(rawEventResponse, &rawResponse)
	if err != nil {
		return nil, err
	}

	events := []map[string]json.RawMessage{}
	err = json.Unmarshal(rawResponse["events"], &events)
	if err != nil {
		return nil, err
	}

	messages := []EventMessage{}
	for _, event := range events {
		// if the event is a heartbeat, return a special error
		if string(event["type"]) == `"heartbeat"` {
			return nil, HeartbeatError
		}
		var msg EventMessage
		err = json.Unmarshal(event["message"], &msg)
		// TODO? should this check be here
		if err != nil {
			return nil, err
		}
		msg.Queue = q
		messages = append(messages, msg)
	}

	return messages, nil
}
