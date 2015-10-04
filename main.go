package main

import (
	"github.com/natural/missmolly/log"
	"github.com/natural/missmolly/server"
)

func main() {
	srv, err := server.NewFromFile("missmolly.conf")
	if err != nil {
		log.Error("main", "error", err)
		panic(err)
	}
	log.Fatal(srv.ListenAndServe())
}
