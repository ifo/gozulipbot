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
	emailAddress, apiKey, err := gzb.GetConfigFromFlags()
	if err != nil {
		log.Fatalln(err)
	}

	bot := gzb.MakeBot(emailAddress, apiKey, []string{})

	streams := getStreamListNames(bot)

	subscribeResp := subscribeToStreams(bot, streams)

	fmt.Println(subscribeResp.String())
}

func getStreamListNames(bot gzb.Bot) []string {
	list, err := bot.GetStreamNameList()
	if err != nil {
		log.Fatal(err)
	}

	return list
}

func subscribeToStreams(bot gzb.Bot, streams []string) bytes.Buffer {
	resp, err := bot.SubscribeToStreams(streams)
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
