package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	var (
		emailAddress = flag.String("email", "", "bot email address")
		apiKey       = flag.String("apikey", "", "bot api key")
	)
	flag.Parse()

	bot := MakeBot(*emailAddress, *apiKey, []string{"test-bot"})

	resp, err := bot.SendStreamMessage("test-bot", "test-go-bot",
		"continued progress")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}
