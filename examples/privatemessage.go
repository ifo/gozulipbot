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

	resp, err := bot.SendPrivateMessage("person@example.com",
		"this message is private")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var toPrint bytes.Buffer

	err = json.Indent(&toPrint, body, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(toPrint.String())
}
