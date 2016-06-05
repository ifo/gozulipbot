package gozulipbot

import (
	"net/http"
	"testing"
)

func TestBot_Init(t *testing.T) {
	bot := Bot{}
	bot.Init()

	if bot.client == nil {
		t.Error("expected bot to have client")
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
		Email:   "testbot@example.com",
		APIKey:  "apikey",
		Streams: []string{"stream a", "test bots"},
		client:  &testClient{},
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
