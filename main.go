package main

import (
	"net/http"
	"remote-storage/router"
)

func main() {
	var route = router.NewHttpRouter()
	route.AddHandler("/test", f)
	route.Listen()
}

func f(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(200)
	w.Write([]byte("Hello, world!"))
}
