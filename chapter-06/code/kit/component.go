package kit

import (
	"encoding/json"
	"expvar"
	"fmt"
	"log"
	"runtime"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nuid"
)

// Component is contains reusable logic related to handling
// of the connection to NATS in the system.
type Component struct {
	// cmu is the lock from the component.
	cmu sync.Mutex

	// id is a unique identifier used for this component.
	id string

	// nc is the connection to NATS.
	nc *nats.Conn

	// kind is the type of component.
	kind string
}

// NewComponent creates a
func NewComponent(kind string) *Component {
	id := nuid.Next()
	return &Component{
		id:   id,
		kind: kind,
	}
}

// SetupConnectionToNATS connects to NATS and registers the event
// callbacks and makes it available for discovery requests as well.
func (c *Component) SetupConnectionToNATS(servers string, options ...nats.Option) error {
	// Label the connection with the kind and id from component.
	options = append(options, nats.Name(c.Name()))

	c.cmu.Lock()
	defer c.cmu.Unlock()

	// Connect to NATS with customized options.
	nc, err := nats.Connect(servers, options...)
	if err != nil {
		return err
	}
	c.nc = nc

	// Setup NATS event callbacks
	//
	// Handle protocol errors and slow consumer cases.
	nc.SetErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
		log.Printf("NATS error: %s\n", err)
	})
	nc.SetReconnectHandler(func(_ *nats.Conn) {
		log.Println("Reconnected to NATS!")
	})
	nc.SetDisconnectHandler(func(_ *nats.Conn) {
		log.Println("Disconnected from NATS!")
	})
	nc.SetClosedHandler(func(_ *nats.Conn) {
		panic("Connection to NATS is closed!")
	})

	// Register component so that it is available for discovery requests.
	_, err = c.nc.Subscribe("_NATS_RIDER.discovery", func(m *nats.Msg) {
		// Reply back directly with own name if requested.
		if m.Reply != "" {
			nc.Publish(m.Reply, []byte(c.ID()))
		} else {
			log.Println("[Discovery] No Reply inbox, skipping...")
		}
	})

	// Register component so that it is available for direct status requests.
	// e.g. _NATS_RIDER.api-server.:id.status
	statusSubject := fmt.Sprintf("_NATS_RIDER.%s.status", c.id)
	_, err = c.nc.Subscribe(statusSubject, func(m *nats.Msg) {
		if m.Reply != "" {
			log.Println("[Status] Replying with status...")
			statsz := struct {
				Kind string           `json:"kind"`
				ID   string           `json:"id"`
				Cmd  []string         `json:"cmdline"`
				Mem  runtime.MemStats `json:"memstats"`
			}{
				Kind: c.kind,
				ID:   c.id,
				Cmd:  expvar.Get("cmdline").(expvar.Func)().([]string),
				Mem:  expvar.Get("memstats").(expvar.Func)().(runtime.MemStats),
			}
			result, err := json.Marshal(statsz)
			if err != nil {
				log.Printf("Error: %s\n", err)
				return
			}
			nc.Publish(m.Reply, result)
		} else {
			log.Println("[Status] No Reply inbox, skipping...")
		}
	})

	return err
}

// NATS returns the current NATS connection.
func (c *Component) NATS() *nats.Conn {
	c.cmu.Lock()
	defer c.cmu.Unlock()
	return c.nc
}

// ID returns the ID from the component.
func (c *Component) ID() string {
	c.cmu.Lock()
	defer c.cmu.Unlock()
	return c.id
}

// Name is the label used to identify the NATS connection.
func (c *Component) Name() string {
	c.cmu.Lock()
	defer c.cmu.Unlock()
	return fmt.Sprintf("%s:%s", c.kind, c.id)
}

// Shutdown makes the component go away.
func (c *Component) Shutdown() error {
	c.NATS().Close()
	return nil
}
