package main

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"

	"github.com/sean9999/metropolis/event"
)

type EventHandler struct {
	Metropolis
}

// handle /events
func (m EventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	bus := m.Bus

	switch r.Method {
	case "GET":

		//	SSE
		w.Header().Set("Access-Control-Allow-Origin", "https://rewinder.lcl.host:44324")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		fmt.Fprintf(w, "\n%s\n\n", ": and we begin")
		for ev := range bus.Queue {
			fmt.Fprint(w, ev.Serialize())
			w.(http.Flusher).Flush()
		}

	case "POST":

		//	JSON
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
		e := bus.NewEvent(*payload)
		w.Write([]byte(e.Serialize()))
		bus.Enqueue(e)

	case "PUT":

		//	PEM
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
		e := bus.NewEvent(payload)
		w.Write([]byte(e.Serialize()))
		bus.Enqueue(e)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
	}
	r.Body.Close()
}
