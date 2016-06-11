package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

	q, err := bot.RegisterAt()
	if err != nil {
		log.Fatal(err)
	}

	evtResp := getEventsFromQueue(q)

	messages, err := gzb.ParseEventMessages(evtResp.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	// Respond with "hi" to all @ messages
	for _, m := range messages {
		resp, err := bot.Respond(m, "hi")
		if err != nil {
			log.Println(err)
		} else {
			defer resp.Body.Close()
			printResponse(resp.Body)
		}
	}
}

func getEventsFromQueue(q *gzb.Queue) bytes.Buffer {
	resp, err := q.GetEvents()
	if err != nil {
		log.Fatal("get events from queue error: ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("get events from queue error 2: ", err)
	}

	var out bytes.Buffer
	err = json.Indent(&out, body, "", "  ")
	if err != nil {
		log.Fatal("get events from queue error 3: ", err)
	}

	return out
}

func printResponse(r io.ReadCloser) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println(err)
	}
	var toPrint bytes.Buffer
	err = json.Indent(&toPrint, body, "", "  ")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(toPrint.String())
}
