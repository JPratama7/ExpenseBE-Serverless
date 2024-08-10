package main

import (
	"crud/route"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func main() {
	functions.HTTP("Serve", route.Router.ServeHTTP)
}
