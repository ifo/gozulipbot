package gozulipbot

import (
	"net/http"
	"testing"
)

func TestMakeBot(t *testing.T) {
	cases := map[string][]string{
		"1": {"email@example.com", "secretkey", "one-stream"},
	}

	for k, c := range cases {
		b := MakeBot(c[0], c[1], []string{c[2]})
		if b.EmailAddress != c[0] {
			t.Errorf("case %s, actual %s, expected %s", k, b.EmailAddress, c[0])
		}
		if b.ApiKey != c[1] {
			t.Errorf("case %s, actual %s, expected %s", k, b.ApiKey, c[1])
		}
		if b.Streams[0] != c[2] {
			t.Errorf("case %s, actual %s, expected %s", k, b.Streams[0], c[2])
		}
	}
}

func TestMessage(t *testing.T) {
	t.Skip()
}

func TestPrivateMessage(t *testing.T) {
	t.Skip()
}

func getTestBot() *Bot {
	return &Bot{
		EmailAddress: "testbot@example.com",
		ApiKey:       "apikey",
		Streams:      []string{"stream a", "test bots"},
		client:       &testClient{},
	}
}

type testClient struct {
	Request  *http.Request
	Response *http.Response
}

func (t *testClient) Do(r *http.Request) (*http.Response, error) {
	t.Request = r
	return t.Response, nil
}
