package gozulipbot

import (
	"errors"
	"flag"
)

func GetConfigFromFlags() (string, string, error) {
	var (
		emailAddress = flag.String("email", "", "bot email address")
		apiKey       = flag.String("apikey", "", "bot api key")
	)
	flag.Parse()

	if *emailAddress == "" {
		return "", "", errors.New("--email required, but wasn't set")
	}
	if *apiKey == "" {
		return "", "", errors.New("--apikey required, but wasn't set")
	}

	return *emailAddress, *apiKey, nil
}
