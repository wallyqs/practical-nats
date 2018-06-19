package main


import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/go-nats"
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
	nc.Subscribe("long.request", func(m *nats.Msg) {
		log.Println("[Processing]", string(m.Data))
		var req RequestWithKeepAlive
		err := json.Unmarshal(m.Data, &req)
		if err != nil {
		        log.Printf("Error: %s", err)
			nc.Publish(m.Reply, []byte("error!"))
			return
		}
		log.Printf("[Heartbeats] %+v", req)

		// Do the work
		time.Sleep(30 * time.Second)
		nc.Publish(m.Reply, []byte("done!"))
	})

	log.Println("[Started]")
	select {}
}
