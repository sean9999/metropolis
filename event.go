package main

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"

	"github.com/sean9999/metropolis/event"
)

type Ingester struct {
	Bus *event.Bus
}

// processJSON takes in JSON payload and outputs text encoded Event
func (s Ingester) processJSON(w http.ResponseWriter, r *http.Request) {
	m := s.Bus
	rawBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(401)
		fmt.Fprint(w, "could not read body")
		return
	}
	r.Body.Close()
	payload := new(event.Payload)
	err = json.Unmarshal(rawBytes, payload)
	if err != nil {
		w.WriteHeader(401)
		fmt.Fprint(w, "could not unmarshal JSON")
		return
	}
	e := m.NewEvent(*payload)
	w.Write([]byte(e.Serialize()))
	n := m.Enqueue(e)
	fmt.Println(n)
}

// processPEM takes in PEM-encoded text and outputs an Event
func (s Ingester) processPEM(w http.ResponseWriter, r *http.Request) {
	m := s.Bus
	rawBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(401)
		fmt.Fprint(w, "could not read body")
		return
	}
	r.Body.Close()
	pem, rest := pem.Decode(rawBytes)
	payload := event.Payload{
		"type":    pem.Type,
		"headers": pem.Headers,
		"bytes":   pem.Bytes,
		"rest":    rest,
	}
	e := m.NewEvent(payload)
	w.Write([]byte(e.Serialize()))
	m.Enqueue(e)
}

func (s Ingester) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "POST,PUT,GET")
	switch r.Method {
	case "POST":
		s.processJSON(w, r)
	case "PUT":
		s.processPEM(w, r)
	}
}

type Emitter struct {
	Bus *event.Bus
}

func (em Emitter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bus := em.Bus
	w.Header().Set("Access-Control-Allow-Origin", "https://rewinder.lcl.host:44324")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	fmt.Fprintf(w, "\n%s\n\n", ": and we begin")

	//	run loop
	for {
		//	if there is nothing in the queue, do nothing & go again
		if len(bus.Queue) > 0 {
			bus.RLock()
			//	if there is something in the queue AND we're allowed to emit, emit
			//	otherwise, leave it in the queue
			if bus.SSEEmit {
				ev := <-bus.Queue
				fmt.Fprint(w, ev.Serialize())
				w.(http.Flusher).Flush()
			}
			bus.RUnlock()
		}
	}
}

type EventHandler struct {
	Metropolis
}

func (m EventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		Emitter{m.Bus}.ServeHTTP(w, r)
	case "POST":
	case "PUT":
		//m.Bus.Ingester.ServeHTTP(w, r)
		Ingester{m.Bus}.ServeHTTP(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
	}
	r.Body.Close()
}
