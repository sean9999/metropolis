package main

import (
	"crypto/tls"
	"net/http"
	"os"

	"github.com/rs/cors"
	"github.com/sean9999/metropolis/web"
)

func main() {

	//	singleton
	met, err := NewMetropolis()
	if err != nil {
		panic(err)
	}

	//	routing
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		web.StaticServer{}.ServeHTTP(w, r)
	})
	mux.Handle("GET /index.html", web.StaticServer{})
	mux.Handle("GET /js.js", web.StaticServer{})
	mux.Handle("GET /mod.js", web.StaticServer{})
	mux.Handle("GET /graph.js", web.StaticServer{})
	mux.Handle("GET /utils.js", web.StaticServer{})
	mux.Handle("GET /favicon.ico", web.StaticServer{})
	mux.Handle("GET /css.css", web.StaticServer{})
	mux.Handle("/graph/*", GraphHandler{met})
	mux.Handle("/events", EventHandler{met})

	//	cors middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://rewinder.lcl.host:44324"},
		AllowCredentials: true,
		Debug:            false,
		AllowedMethods:   []string{"GET", "PUT", "POST", "OPTIONS"},
	})
	corehan := c.Handler(mux)

	//	server
	ln, _ := tls.Listen("tcp", ":"+os.Getenv("HTTPS_PORT"), Anchor())
	http.Serve(ln, corehan)

}
