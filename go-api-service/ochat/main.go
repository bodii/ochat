package main

import (
	"fmt"
	"net/http"
	"ochat/bootstrap"
	"ochat/router"
)

// main func
func main() {
	bootstrap.Init()

	router.Init()
	router.WebsocketInit()

	systemConf := bootstrap.SystemConf
	servConf := systemConf.Serv
	http.ListenAndServe(
		fmt.Sprintf("%s:%d", servConf.Host, servConf.Port),
		nil)
}
