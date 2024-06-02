package main

import (
	"encoding/pem"
	"fmt"
	"io"
	"net/http"

	"github.com/sean9999/metropolis/event"
)

type GraphHandler struct {
	Metropolis
}

func (m GraphHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	switch r.URL.Path {
	case "/graph/me":
		//	show me as a peer (public info only)
		w.Header().Set("Content-Type", "application/json")
		j, _ := m.Network.Me.AsPeer().MarshalJSON()
		fmt.Fprintln(w, string(j))
	case "/graph/assert":

		switch r.Method {
		case "PUT":
			bus := m.Bus
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
		}

	default:
		fmt.Fprintf(w, "the method is %s\n", r.Method)
		fmt.Fprintf(w, "the URL is %s\n", r.URL)
	}

}
