package main

import (
	"log"
	"time"

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

	q, err := bot.RegisterAt()
	if err != nil {
		log.Fatal(err)
	}

	msgs, stopFunc := q.EventsChan()

	// stop after a minute
	go func() {
		time.Sleep(1 * time.Minute)
		stopFunc()
	}()

	for m := range msgs {
		log.Println("message received")
		m.Queue.Bot.Respond(m, "hi forever!")
	}
}
