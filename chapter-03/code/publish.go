package main


import (
        "log"
        "github.com/nats-io/go-nats"
)

func main(){
        nc, err := nats.Connect("nats://127.0.0.1:4222")
        if err != nil {
                log.Fatalf("Error: %s", err)
        }
        nc.Publish("hello", []byte("world"))
}
