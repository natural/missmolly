package main

import (
	"log"
	"net/http"
	"time"

	"github.com/natural/missmolly/server"
)

func main() {
	app := server.NewFromFile("missmolly.conf")
	srv := &http.Server{
		Addr:           ":7373",
		Handler:        app, //app per server, or...
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(srv.ListenAndServe())
}
