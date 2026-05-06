package main

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type JsonRPC2 struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Result  interface{} `json:"result,omitempty"`
	ID      *int        `json:"id,omitempty"`
}

type SubscribeParams struct {
	Channel string `json:"channel"`
}

func main() {
	u := url.URL{Scheme: "wss", Host: "ws.lightstream.bitflyer.com", Path: "/json-rpc"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	subscribe := &JsonRPC2{
		Version: "2.0",
		Method:  "subscribe",
		Params:  &SubscribeParams{Channel: "lightning_ticker_BTC_JPY"},
	}
	if err := c.WriteJSON(subscribe); err != nil {
		log.Fatal("subscribe:", err)
	}

	for {
		message := new(JsonRPC2)
		if err := c.ReadJSON(message); err != nil {
			log.Println("read:", err)
			return
		}
		if message.Method == "channelMessage" {
			log.Println(message.Params)
		}
	}
}