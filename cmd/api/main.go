package main

import (
	"crud/route"
	"log"
	"net/http"
	"time"
)

func main() {

	srv := &http.Server{
		Handler: route.Router,
		Addr:    "127.0.0.1:8000",

		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Fatalln(srv.ListenAndServe())

}
