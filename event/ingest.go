package event

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
)

type Ingester struct {
	Store *MemoryStore
}

// processJSON takes in JSON payload and outputs text encoded Event
func (s Ingester) processJSON(w http.ResponseWriter, r *http.Request) {
	m := s.Store
	rawBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(401)
		fmt.Fprint(w, "could not read body")
		return
	}
	r.Body.Close()
	payload := new(Payload)
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
	m := s.Store
	rawBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(401)
		fmt.Fprint(w, "could not read body")
		return
	}
	r.Body.Close()

	pem, rest := pem.Decode(rawBytes)

	payload := Payload{
		"type":    pem.Type,
		"headers": pem.Headers,
		"bytes":   pem.Bytes,
		"rest":    rest,
	}
	e := m.NewEvent(payload)
	w.Write([]byte(e.Serialize()))
	n := m.Enqueue(e)

	fmt.Println(n)

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
