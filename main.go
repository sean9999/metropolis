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
	"github.com/sean9999/metropolis/event"
	"github.com/sean9999/metropolis/web"
)

func main() {
	//	Event Bus

	met, err := NewMetropolis()
	if err != nil {
		panic(err)
	}

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
				met.Bus.Lock()
				met.Bus.SSEEmit = false
				met.Bus.Unlock()
			case syscall.SIGUSR2:
				//	continue. do emit
				met.Bus.Lock()
				met.Bus.SSEEmit = true
				met.Bus.Unlock()
			}
		}
	}()

	//	send 8 events off the top
	go func() {
		for n := 0; n < 8; n++ {
			time.Sleep(time.Millisecond * 250)
			e := met.Bus.NewEvent(event.Payload{})
			e.Data["msg"] = fmt.Sprintf("I am event Id %d", e.Id)
			met.Bus.Enqueue(e)
		}
	}()

	mux := http.NewServeMux()

	mux.Handle("/graph/*", GraphHandler{met})
	mux.Handle("/events", EventHandler{met})

	//	static
	mux.Handle("GET /index.html", web.StaticServer{})
	mux.Handle("GET /js.js", web.StaticServer{})
	mux.Handle("GET /css.css", web.StaticServer{})

	os.WriteFile("rewinder.pid", []byte(fmt.Sprintf("%d", os.Getpid())), 0640)

	fmt.Printf("pid: %d\n", os.Getpid())

	ln, _ := tls.Listen("tcp", ":"+os.Getenv("HTTPS_PORT"), Anchor())

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://rewinder.lcl.host:44324"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug:          false,
		AllowedMethods: []string{"GET", "PUT", "POST", "OPTIONS"},
	})

	// Insert the middleware
	corehan := c.Handler(mux)

	// Start the https server
	http.Serve(ln, corehan)

}
