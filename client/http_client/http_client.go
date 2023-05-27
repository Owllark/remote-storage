package http_client

import (
	"log"
	"net/http"
	"time"
)

type Client interface {
	Listen()
	AddHandler()
}

func NewHttpRouter() *HttpRouter {
	h := tServerHandler{
		mux: make(map[string]func(http.ResponseWriter, *http.Request)),
	}
	r := HttpRouter{
		serverHandler: h,
	}
	return &r
}

type tServerHandler struct {
	mux map[string]func(http.ResponseWriter, *http.Request)
}

func (h *tServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handlerFunc, ok := h.mux[r.URL.String()]; ok {
		handlerFunc(w, r)
		return
	}
}

type HttpRouter struct {
	server        http.Server
	serverHandler tServerHandler
}

func (r *HttpRouter) Listen() {
	r.server = http.Server{
		Addr:        ":8080",
		Handler:     &r.serverHandler,
		ReadTimeout: 5 * time.Second,
	}
	err := r.server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func (r *HttpRouter) AddHandler(url string, f func(http.ResponseWriter, *http.Request)) {
	r.serverHandler.mux[url] = f
}
