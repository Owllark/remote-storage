package main

import "net/http"

func f(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(200)
	w.Write([]byte("Hello, world!"))
}

func (w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(200)
	w.Write([]byte("Hello, world!"))
}
