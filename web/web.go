package web

import (
	"fmt"
	"net/http"
)

type StaticServer struct{}

func (h StaticServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "" {
		r.URL.Path = "index.html"
	}
	http.ServeFile(w, r, fmt.Sprintf("web/%s", r.URL.Path))
}

//	static
// mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "web/index.html")
// })
// mux.HandleFunc("GET /js.js", func(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "web/js.js")
// })
// mux.HandleFunc("GET /css.css", func(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "web/css.css")
// })
