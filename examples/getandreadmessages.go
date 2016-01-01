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

	regResp := registerEvents(bot)
	regRespJson := map[string]interface{}{}

	err = json.Unmarshal(regResp.Bytes(), &regRespJson)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(regRespJson)

	queueID := regRespJson["queue_id"].(string)
	lastEventID := int(regRespJson["last_event_id"].(float64))
	maxMsgID := int(regRespJson["max_message_id"].(float64))
	if lastEventID < maxMsgID {
		lastEventID = maxMsgID
	}

	evtResp := getEventsFromQueue(bot, queueID, lastEventID)

	fmt.Println(evtResp.String())
	// uncomment in case response is huuuuuuuge
	//ioutil.WriteFile("events-response.json", evtResp.Bytes(), 0644)
}

func registerEvents(bot gzb.Bot) bytes.Buffer {
	resp, err := bot.RegisterEvents()
	if err != nil {
		log.Fatal("register events error 1: ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("register events error 2: ", err)
	}

	var out bytes.Buffer
	err = json.Indent(&out, body, "", "  ")
	if err != nil {
		log.Fatal("register events error 3: ", err)
	}

	return out
}

func getEventsFromQueue(bot gzb.Bot, queueID string,
	lastMessageID int) bytes.Buffer {
	resp, err := bot.GetEventsFromQueue(queueID, lastMessageID)
	if err != nil {
		log.Fatal("get events from queue error: ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("get events from queue error 2: ", err)
	}

	// uncomment in case of 404, which causes json.Indent to fail
	//out := bytes.NewBuffer(body)
	//return *out

	var out bytes.Buffer
	err = json.Indent(&out, body, "", "  ")
	if err != nil {
		log.Fatal("get events from queue error 3: ", err)
	}

	return out
}
