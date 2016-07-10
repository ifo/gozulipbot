package gozulipbot

import (
	"net/http"
	"testing"
)

func TestBot_Init(t *testing.T) {
	bot := Bot{}
	bot.Init()

	if bot.Client == nil {
		t.Error("expected bot to have client")
	}
}

func TestBot_GetStreamList(t *testing.T) {
	t.Skip()
}

func TestBot_GetStreams(t *testing.T) {
	t.Skip()
}

func TestBot_Subscribe(t *testing.T) {
	t.Skip()
}

func TestBot_Unsubscribe(t *testing.T) {
	t.Skip()
}

func TestBot_ListSubscriptions(t *testing.T) {
	t.Skip()
}

func TestBot_RegisterEvents(t *testing.T) {
	t.Skip()
}

func TestBot_RegisterAll(t *testing.T) {
	t.Skip()
}

func TestBot_RegisterAt(t *testing.T) {
	t.Skip()
}

func TestBot_RegisterPrivate(t *testing.T) {
	t.Skip()
}

func TestBot_RegisterSubscriptions(t *testing.T) {
	t.Skip()
}

func TestBot_RawRegisterEvents(t *testing.T) {
	t.Skip()
}

// ensure constructRequest adds a JSON header and uses basic auth
func TestBot_constructRequest(t *testing.T) {
	bot := getTestBot()
	type Case struct {
		Method   string
		Endpoint string
		Body     string
		ReqBody  string
		Error    error
	}

	JSONHeader := "application/x-www-form-urlencoded"

	cases := map[string]Case{
		"1": Case{"GET", "endpoint", "", "", nil},
	}

	for num, c := range cases {
		req, err := bot.constructRequest(c.Method, c.Endpoint, c.Body)
		if err != nil {
			t.Fatalf("got %q, expected nil, case %q", err, num)
		}

		header := req.Header.Get("Content-Type")
		if string(header) != JSONHeader {
			t.Errorf("got %q, expected %q, case %q", header, JSONHeader, num)
		}

		email, key, ok := req.BasicAuth()
		if !ok || email != bot.Email || key != bot.APIKey {
			t.Errorf("got %t, expected true, case %q", ok, num)
			t.Errorf("got %q, expected %q, case %q", email, bot.Email, num)
			t.Errorf("got %q, expected %q, case %q", key, bot.APIKey, num)
		}
	}
}

func getTestBot() *Bot {
	return &Bot{
		Email:   "testbot@example.com",
		APIKey:  "apikey",
		Streams: []string{"stream a", "test bots"},
		Client:  &testClient{},
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
