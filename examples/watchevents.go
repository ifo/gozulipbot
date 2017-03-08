package main

import (
	"log"
	"time"

	gzb "github.com/ifo/gozulipbot"
)

func main() {
	bot := gzb.Bot{}
	err := bot.GetConfigFromFlags()
	if err != nil {
		log.Fatalln(err)
	}
	bot.Init()

	q, err := bot.RegisterAt()
	if err != nil {
		log.Fatal(err)
	}

	stopFunc := q.EventsCallback(respondToMessage)

	time.Sleep(1 * time.Minute)
	stopFunc()
}

func respondToMessage(em gzb.EventMessage, err error) {
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("message received")

	em.Queue.Bot.Respond(em, "hi forever!")
}
