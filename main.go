package main

import (
	"log"

	"github.com/natural/missmolly/server"
)

func main() {
	srv, err := server.NewFromFile("missmolly.conf")
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(srv.ListenAndServe())
}
