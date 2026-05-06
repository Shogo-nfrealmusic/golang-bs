package main

import (
	"fmt"
	"log"
	"os"

	pubnub "github.com/pubnub/go/v8"
)

func main() {
	subKey := os.Getenv("PUBNUB_SUBSCRIBE_KEY")
	if subKey == "" {
		log.Fatal("PUBNUB_SUBSCRIBE_KEY is required")
	}

	config := pubnub.NewConfig("golang-bs-user")
	config.SubscribeKey = subKey

	pn := pubnub.NewPubNub(config)
	fmt.Println("PubNub initialized:", pn != nil)
}