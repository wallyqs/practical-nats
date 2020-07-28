package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	nc.Publish("_SYS.hi", []byte("hi"))
	nc.Flush()
	time.Sleep(1 * time.Second)
	log.Printf("Last Error: %s", nc.LastError())
}
