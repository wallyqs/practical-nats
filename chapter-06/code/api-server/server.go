package apiserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/nats-io/nuid"
	"github.com/wallyqs/practical-nats/chapter-06/code/kit"
)

const (
	Version = "0.1.0"
)

// Server is a component.
type Server struct {
	*kit.Component
}

// HandleRides processes requests to find available drivers in an area.
func (s *Server) HandleRides(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var request *kit.DriverAgentRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Tag the request with an ID for tracing in the logs.
	request.RequestID = nuid.Next()
	req, err := json.Marshal(request)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	nc := s.NATS()

	// Find a driver available to help with the request.
	log.Printf("requestID:%s - Finding available driver for request: %s\n", request.RequestID, string(body))
	msg, err := nc.Request("drivers.find", req, 5*time.Second)
	if err != nil {
		log.Printf("requestID:%s - Gave up finding available driver for request\n", request.RequestID)
		http.Error(w, "Request timeout", http.StatusRequestTimeout)
		return
	}
	log.Printf("requestID:%s - Response: %s\n", request.RequestID, string(msg.Data))

	var resp *kit.DriverAgentResponse
	err = json.Unmarshal(msg.Data, &resp)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if resp.Error != "" {
		http.Error(w, resp.Error, http.StatusServiceUnavailable)
		return
	}

	log.Printf("requestID:%s - Driver with ID %s is available to handle the request", request.RequestID, resp.ID)
	fmt.Fprintf(w, string(msg.Data))
}

// ListenAndServe takes the network address and port that
// the HTTP server should bind to and starts it.
func (s *Server) ListenAndServe(addr string) error {
	mux := http.NewServeMux()

	// GET /
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// See: https://golang.org/pkg/net/http/#ServeMux.Handle
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, fmt.Sprintf("NATS Rider API Server v%s\n", Version))
	})

	// POST /rides
	mux.HandleFunc("/rides", s.HandleRides)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	srv := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go srv.Serve(l)

	return nil
}
