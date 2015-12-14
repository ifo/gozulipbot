package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	var (
		emailAddress = flag.String("email", "", "bot email address")
		apiKey       = flag.String("apikey", "", "bot api key")
	)
	flag.Parse()

	bot := Bot{
		EmailAddress: *emailAddress,
		ApiKey:       *apiKey,
		Streams:      []string{"test-bot"},
	}

	c := http.Client{}

	v := url.Values{}
	v.Set("type", "stream")
	v.Set("to", "test-bot")
	v.Set("subject", "test-go-bot")
	v.Set("content", "okay now this works")

	req, err := http.NewRequest("POST", "https://api.zulip.com/v1/messages",
		strings.NewReader(v.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(bot.EmailAddress, bot.ApiKey)

	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}
