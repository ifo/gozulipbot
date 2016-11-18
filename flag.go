package gozulipbot

import (
	"errors"
	"flag"
)

func (b *Bot) GetConfigFromFlags() error {
	var (
		apiKey = flag.String("apikey", "", "bot api key")
		apiURL = flag.String("apiurl", "", "url of zulip server")
		email  = flag.String("email", "", "bot email address")
	)
	flag.Parse()

	if *apiKey == "" {
		return errors.New("--apikey is required")
	}
	if *apiURL == "" {
		return errors.New("--apiurl is required")
	}
	if *email == "" {
		return errors.New("--email is required")
	}
	b.APIKey = *apiKey
	b.APIURL = *apiURL
	b.Email = *email
	return nil
}
