package main

import (
	"fmt"
	"log"

	gzb "github.com/ifo/gozulipbot"
)

func main() {
	bot := gzb.Bot{}
	err := bot.GetConfigFromFlags()
	if err != nil {
		log.Fatalln(err)
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
