package gozulipbot

import (
	"encoding/json"
)

type EventMessage struct {
	AvatarURL        string           `json:"avatar_url"`
	Client           string           `json:"client"`
	Content          string           `json:"content"`
	ContentType      string           `json:"content_type"`
	DisplayRecipient DisplayRecipient `json:"display_recipient"`
	GravatarHash     string           `json:"gravatar_hash"`
	ID               int              `json:"id"`
	RecipientID      int              `json:"recipient_id"`
	SenderDomain     string           `json:"sender_domain"`
	SenderEmail      string           `json:"sender_email"`
	SenderFullName   string           `json:"sender_full_name"`
	SenderID         int              `json:"sender_id"`
	SenderShortName  string           `json:"sender_short_name"`
	Subject          string           `json:"subject"`
	SubjectLinks     []interface{}    `json:"subject_links"`
	Timestamp        int              `json:"timestamp"`
	Type             string           `json:"type"`
}

type DisplayRecipient struct {
	Users []User `json:"users,omitempty"`
	Topic string `json:"topic,omitempty"`
}

type User struct {
	Domain        string `json:"domain"`
	Email         string `json:"email"`
	FullName      string `json:"full_name"`
	ID            int    `json:"id"`
	IsMirrorDummy bool   `json:"is_mirror_dummy"`
	ShortName     string `json:"short_name"`
}

func (d *DisplayRecipient) UnmarshalJSON(b []byte) (err error) {
	topic, users := "", make([]User, 1)
	if err = json.Unmarshal(b, &topic); err == nil {
		d.Topic = topic
		return
	}
	if err = json.Unmarshal(b, &users); err == nil {
		d.Users = users
		return
	}
	return
}

func ParseEventMessages(rawEventResponse []byte) ([]EventMessage, error) {
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
		var msg EventMessage
		err = json.Unmarshal(event["message"], &msg)
		// TODO: determine if this check should be here
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
