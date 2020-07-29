package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type RequestWithKeepAlive struct {
	HeartbeatsInbox string `json:"hb_inbox"`
	Data            []byte `json:"data"`
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	hbInbox := nats.NewInbox()
	req := &RequestWithKeepAlive{
		HeartbeatsInbox: hbInbox,
		Data:            []byte("hello world"),
	}
	payload, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, err := nc.RequestWithContext(ctx, "long.request", payload)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[Response]", string(response.Data))
}
