package main


import (
	"log"

	"github.com/nats-io/go-nats"
)

func main() {
	nc, err := nats.Connect("nats://127.0.0.1:4222",
		nats.DisconnectHandler(func(nc *nats.Conn) {
			log.Printf("Got disconnected!\n")
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("Got reconnected to %v!\n", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Printf("Connection closed. Reason: %v\n", nc.LastError())
		}),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			log.Printf("Error: %s\n", err)
			if err == nats.ErrSlowConsumer {
				log.Printf("Removing subscription on %q\n", sub.Subject)
				sub.Unsubscribe()
			}
		}),
	)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	sub, _ := nc.Subscribe("foo", func(_ *nats.Msg) {
		// Heavy processing
		log.Println("Start processing 'foo' message")
		for i := 0; i < 10000000000; i++ {
		}
		log.Println("Done processing 'foo' message")
	})
	sub.SetPendingLimits(8192, 8192*1024)

	nc.Subscribe("bar", func(_ *nats.Msg) {
		// Not heavy processing
	})
	select {}
}
