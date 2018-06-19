package main


import (
	"encoding/json"
	"log"
	"runtime"

	"github.com/nats-io/go-nats"
)

func main() {
	nc, err := nats.Connect("nats://127.0.0.1:4222",
		nats.UserInfo("foo", "secret"),
	)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	nc.Subscribe("greeting", func(m *nats.Msg) {
		log.Printf("[Received] %s", string(m.Data))
	})

	payload := struct {
		RequestID string
		Data      []byte
	}{
		RequestID: "1234-5678-90",
		Data:      []byte("encoded data"),
	}
	msg, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	nc.Publish("greeting", msg)
	runtime.Goexit()
}
