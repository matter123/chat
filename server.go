package main

import (
	"github.com/matter123/chat/chatserver"
	"github.com/matter123/chat/config"
	"log"
	"net/http"
	"time"
)

func main() {
	server := &http.Server{
		Addr:           ":" + config.Settings().Port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	http.Handle("/", http.FileServer(http.Dir("web_base")))
	http.HandleFunc("/login", chatserver.LoginHandle)
	http.HandleFunc("/signup", chatserver.SignupHandle)
	if config.Settings().SSLSettings == nil {
		log.Fatal(server.ListenAndServe())
	} else {
		log.Fatal(server.ListenAndServeTLS(config.Settings().SSLSettings.Certificate,
			config.Settings().SSLSettings.Key))
	}
}
