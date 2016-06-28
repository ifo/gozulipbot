package gozulipbot

import (
	"errors"
	"io/ioutil"
	"testing"
)

func TestMessage(t *testing.T) {
	t.Skip()
}

func TestPrivateMessage(t *testing.T) {
	bot := getTestBot()
	type C struct {
		M    Message
		Body string
		E    error
	}
	cases := map[string]C{
		"1": C{M: Message{Stream: "a", Emails: []string{"a@example.com"}, Content: "hey"}, // normal
			Body: "content=hey&to=a%40example.com&type=private", E: nil},
		"2": C{M: Message{Stream: "a", Topic: "a", Emails: []string{"a@example.com"}, Content: "hey"}, // topic is ignored
			Body: "content=hey&to=a%40example.com&type=private", E: nil},
		"3": C{M: Message{Stream: "a", Topic: "a", Emails: []string{"a@example.com", "b@example.com"}, Content: "hey"}, // multiple emails are fine
			Body: "content=hey&to=a%40example.com%2Cb%40example.com&type=private", E: nil},
		"4": C{M: Message{Stream: "a", Content: "hey"}, // no email set
			Body: "", E: errors.New("there must be at least one recipient")},
	}

	for num, c := range cases {
		// ignore response from testClient
		_, err := bot.PrivateMessage(c.M)

		body, _ := ioutil.ReadAll(bot.Client.(*testClient).Request.Body)
		if string(body) != c.Body {
			t.Errorf("got %q, expected %q, case %q", string(body), c.Body, num)
		}

		// no error expected
		if c.E == nil && err != nil {
			t.Errorf("got %q, expected nil, case %q", err, num)
		}
		// error expected, prevent nil.Error() panic
		if c.E != nil && err == nil {
			t.Errorf("got nil, expected %q, case %q", err, c.E, num)
		} else if c.E != nil && err.Error() != c.E.Error() {
			t.Errorf("got %q, expected %q, case %q", err, c.E, num)
		}
	}
}
