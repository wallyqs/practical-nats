package main


import (
	"log"
	"runtime"

	"github.com/nats-io/go-nats"
)

func main() {
	nc1, err := nats.Connect("nats://127.0.0.1:4222", nats.Name("sub"))
	if err != nil {
		log.Fatal(err)
	}
	nc2, err := nats.Connect("nats://127.0.0.1:4222", nats.Name("pub"))
	if err != nil {
		log.Fatal(err)
	}
	nc2.Subscribe("example", func(m *nats.Msg) {
		log.Printf("[Received] %s\n", string(m.Data))
	})

	for i := 0; i < 10; i++ {
		nc1.Publish("example", []byte("hello"))
	}
	nc1.Flush()
	runtime.Goexit()
}
