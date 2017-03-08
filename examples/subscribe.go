package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	streams, err := bot.GetStreams()
	if err != nil {
		log.Fatal(err)
	}

	// print the stream list
	for _, s := range streams {
		fmt.Println(s)
	}

	// subscribe
	subResp := subscribeToStreams(bot, streams)
	fmt.Println(subResp.String())
}

func subscribeToStreams(bot gzb.Bot, streams []string) bytes.Buffer {
	resp, err := bot.Subscribe(streams)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	err = json.Indent(&out, body, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return out
}
