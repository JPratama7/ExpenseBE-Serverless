package main

import (
	"crud/controller"
	"crud/route"
	"log"
	"net/http"
	"time"
)

func main() {

	router := route.NewMuxer()

	router.AddRoute(route.POST, "/register", controller.Register)
	router.AddRoute(route.POST, "/login", controller.Login)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",

		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Fatalln(srv.ListenAndServe())

}
