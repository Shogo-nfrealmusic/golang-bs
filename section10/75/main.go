package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

var DB = map[string]string{
	"User1key": "User1secret",
	"User2key": "User2secret",
}

func Server(apiKey, sign string, data []byte) {
	apiSecret := DB[apiKey]
	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write(data)
	expectedHMAC := hex.EncodeToString(h.Sum(nil))
	if sign != expectedHMAC {
		panic("invalid signature")
	}
	fmt.Println("ok")
}

func main() {
	const apiKey = "User1key"
	const apiSecret = "User1secret"

	data := []byte("data")
	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write(data)
	sign := hex.EncodeToString(h.Sum(nil))
	fmt.Println(sign)

	Server(apiKey, sign, data)
}