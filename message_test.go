package gozulipbot

import (
	"errors"
	"io/ioutil"
	"testing"
)

func TestBot_Message(t *testing.T) {
	t.Skip()
}

func TestBot_PrivateMessage(t *testing.T) {
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

		// Check if error matches the error specified in the case
		switch c.E {
		case nil:
			if err != nil && c.E == nil {
				t.Fatalf("got %q, expected nil, case %q", err, num)
			}

		default:
			if err == nil {
				t.Fatalf("got nil, expected %q, case %q", err, c.E, num)
			}

			if err.Error() != c.E.Error() {
				t.Fatalf("got %q, expected %q, case %q", err, c.E, num)
			}

			// No request was created so the test has been completed and
			// we won't check the request body, as there is none.
			return
		}

		// Check the request body matches our expectation
		body, _ := ioutil.ReadAll(bot.Client.(*testClient).Request.Body)
		if string(body) != c.Body {
			t.Errorf("got %q, expected %q, case %q", string(body), c.Body, num)
		}
	}
}

func TestBot_Respond(t *testing.T) {
	t.Skip()
}

func TestBot_privateResponseList(t *testing.T) {
	t.Skip()
}

func TestBot_constructMessageRequest(t *testing.T) {
	t.Skip()
}
