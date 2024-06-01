package event

import (
	"fmt"
	"net/http"
)

type Emitter struct {
	Store *MemoryStore
}

func (s Emitter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m := s.Store

	fmt.Println("sse connected")

	// w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
	// w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS, GET, DELETE, PUT")
	// w.Header().Add("Access-Control-Allow-Headers", "content-type")
	// w.Header().Add("Access-Control-Max-Age", "86400")
	// w.Header().Add("Access-Control-Allow-Headers", "Authorization, content-type")

	w.Header().Set("Access-Control-Allow-Origin", "https://rewinder.lcl.host:44324")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	fmt.Fprintf(w, "\n%s\n\n", ": and we begin")

	//	run loop
	for {
		//	if there is nothing in the queue, do nothing & go again
		if len(m.Queue) > 0 {
			m.RLock()
			if m.SSEEmit {
				//	if there is something in the queue AND we're allowed to emit, emit
				ev := <-m.Queue
				fmt.Fprint(w, ev.Serialize())
				w.(http.Flusher).Flush()
			}
			m.RUnlock()
		}
	}

}
