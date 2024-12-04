package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sean9999/go-oracle"
	"github.com/sean9999/metropolis/event"
)

type GraphHandler struct {
	Metropolis
}

// handle /graph
func (m GraphHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	switch r.URL.Path {
	case "/graph/me":

		//	show me as a peer (public info only)
		w.Header().Set("Content-Type", "application/json")
		j, _ := m.Network.Me.AsPeer().MarshalJSON()
		fmt.Fprintln(w, string(j))

	case "/graph/verify":

		//	verify assertions made by oracle.Peers
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
			pt := new(oracle.PlainText)
			pt.UnmarshalPEM(rawBytes)
			sender := oracle.NewPeer(nil)
			//	todo: check for nil
			sender.UnmarshalHex([]byte(pt.Headers["pubkey"]))
			ok := m.Network.Me.Verify(pt, sender)
			payload := event.Payload{
				"action":   "peer verification",
				"nickname": sender.Nickname(),
				"pubkey":   pt.Headers["pubkey"],
				"verified": ok,
			}
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				m.Network.Me.AddPeer(sender)
			}
			e := bus.NewEvent(payload)
			jason, _ := json.Marshal(e)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jason)
			bus.Enqueue(e)
		case "GET":
			bus := m.Bus
			pt, err := m.Network.Me.Assert()
			if err != nil {
				w.Write([]byte(err.Error()))
			}
			pem, err := pt.MarshalPEM()
			if err != nil {
				w.Write([]byte(err.Error()))
			}
			payload := event.Payload{
				"action": "host verification",
				"nonce":  pt.Nonce,
			}
			bus.DispatchEvent(payload)
			w.Write(pem)
		}

	default:
		fmt.Fprintf(w, "the method is %s\n", r.Method)
		fmt.Fprintf(w, "the URL is %s\n", r.URL)
	}

}
