package main

import (
	"fmt"
	"net/http"
	"ochat/bootstrap"
	"ochat/controllers"
)

// main func
func main() {
	bootstrap.Init()

	http.HandleFunc("/user/login", controllers.Login)
	http.HandleFunc("/user/signup", controllers.Register)

	servConf := bootstrap.SystemConf.Serv

	http.ListenAndServe(
		fmt.Sprintf("%s:%d", servConf.Host, servConf.Port),
		nil)
}
