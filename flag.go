package gozulipbot

import (
	"errors"
	"flag"
)

func GetConfigFromFlags() (string, string, error) {
	var (
		email  = flag.String("email", "", "bot email address")
		apiKey = flag.String("apikey", "", "bot api key")
	)
	flag.Parse()

	if *email == "" {
		return "", "", errors.New("--email is required")
	}
	if *apiKey == "" {
		return "", "", errors.New("--apikey is required")
	}

	return *email, *apiKey, nil
}
