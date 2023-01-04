package main

import (
	"net/http"
	"ochat/bootstrap"
	"ochat/router"
)

// main func
func main() {
	bootstrap.Init()

	router.Init()
	router.WebsocketInit()

	http.ListenAndServe(bootstrap.HOST_NAME, nil)
}
