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

	bot := gzb.MakeBot(emailAddress, apiKey, []string{})

	regResp := registerEvents(bot)
	queueID, lastEventID := getEventQueueInfo(regResp.Bytes())
	evtResp := getEventsFromQueue(bot, queueID, lastEventID)

	messages, err := gzb.ParseEventMessages(evtResp.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	streamTopic := "example-topic-name"

	// Respond with "hi" to all messages in the specified stream
	for _, m := range messages {
		if m.Subject == streamTopic {
			resp, err := bot.RespondToMessage(m, "hi")
			if err != nil {
				log.Println(err)
			} else {
				defer resp.Body.Close()
				printResponse(resp.Body)
			}
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
	fmt.Println(string(toPrint.Bytes()))
}

func registerEvents(bot gzb.Bot) bytes.Buffer {
	resp, err := bot.RegisterEvents()
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

func getEventQueueInfo(b []byte) (queueID string, lastEventID int) {
	regRespJson := map[string]interface{}{}
	err := json.Unmarshal(b, &regRespJson)
	if err != nil {
		log.Fatalln(err)
	}

	queueID = regRespJson["queue_id"].(string)
	lastEventID = int(regRespJson["last_event_id"].(float64))
	maxMsgID := int(regRespJson["max_message_id"].(float64))
	if lastEventID < maxMsgID {
		lastEventID = maxMsgID
	}
	return
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

	var out bytes.Buffer
	err = json.Indent(&out, body, "", "  ")
	if err != nil {
		log.Fatal("get events from queue error 3: ", err)
	}

	return out
}
