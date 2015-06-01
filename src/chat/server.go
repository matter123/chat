package main

import (
	"net/http"
	"time"
    "log"
)

func main() {
	server := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/login", login)
	http.HandleFunc("/signup", signup)
	log.Fatal(server.ListenAndServe())
}
