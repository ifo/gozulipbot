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

	messages, err := q.GetEvents()
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
