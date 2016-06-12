package main

import (
	"fmt"
	"log"

	gzb "github.com/ifo/gozulipbot"
)

func main() {
	emailAddress, apiKey, err := gzb.GetConfigFromFlags()
	if err != nil {
		log.Fatalln(err)
	}

	bot := gzb.Bot{
		Email:  emailAddress,
		APIKey: apiKey,
	}

	bot.Init()

	q, err := bot.RegisterAll()
	if err != nil {
		log.Fatal(err)
	}

	messages, err := q.GetEvents()
	if err != nil {
		log.Fatal(err)
	}

	// print just the display recipients
	for _, m := range messages {
		fmt.Println(m.DisplayRecipient)
	}

	// print all the messages
	fmt.Println(messages)
}
