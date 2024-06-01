package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/cors"
	"github.com/sean9999/mothership/event"
	"github.com/sean9999/mothership/web"
)

func main() {
	//	Event Bus
	m := event.NewMemoryStore()

	//	listen for signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs,
		syscall.SIGUSR1, // pause
		syscall.SIGUSR2, // continue
	)
	go func() {
		for sig := range sigs {
			switch sig {
			case syscall.SIGUSR1:
				//	pause. so don't emit
				m.Lock()
				m.SSEEmit = false
				m.Unlock()
			case syscall.SIGUSR2:
				//	continue. do emit
				m.Lock()
				m.SSEEmit = true
				m.Unlock()
			}
		}
	}()

	//	send 8 events off the top
	go func() {
		for n := 0; n < 8; n++ {
			time.Sleep(time.Millisecond * 250)
			e := m.NewEvent(event.Payload{})
			e.Data["msg"] = fmt.Sprintf("I am event Id %d", e.Id)
			m.Enqueue(e)
		}
	}()

	// c := cors.New(cors.Options{
	// 	AllowedOrigins: []string{"https://rewinder.lcl.host:44324"},
	// })

	mux := http.NewServeMux()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://rewinder.lcl.host:44324"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug:          false,
		AllowedMethods: []string{"GET", "PUT", "POST", "OPTIONS"},
	})

	// Insert the middleware
	corehan := c.Handler(mux)

	//corehan := cors.Default().Handler(mux)

	mux.Handle("POST /events", m.Ingester) // POST	is for JSON
	mux.Handle("GET /events", m.Emitter)   // GET	is for SSE
	mux.Handle("PUT /events", m.Ingester)  // PUT	is for PEM

	//	static
	mux.Handle("GET /", web.StaticServer{})
	mux.Handle("GET /index.html", web.StaticServer{})
	mux.Handle("GET /js.js", web.StaticServer{})
	mux.Handle("GET /css.css", web.StaticServer{})

	os.WriteFile("rewinder.pid", []byte(fmt.Sprintf("%d", os.Getpid())), 0640)

	fmt.Printf("pid: %d\n", os.Getpid())

	ln, _ := tls.Listen("tcp", ":"+os.Getenv("HTTPS_PORT"), Anchor())

	// Start the https server
	http.Serve(ln, corehan)

}
