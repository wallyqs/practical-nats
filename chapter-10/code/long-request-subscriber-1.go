package main


import (
	"log"
	"time"

	"github.com/nats-io/go-nats"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	nc.Subscribe("very.long.request", func(m *nats.Msg) {
		log.Println("[Processing]", string(m.Data))
		time.Sleep(20 * time.Second)
		nc.Publish(m.Reply, []byte("done!"))
	})

	log.Println("[Started]")
	select {}
}
