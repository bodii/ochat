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

	systemConf := bootstrap.SystemConf

	http.HandleFunc("/user/login", controllers.Login)
	http.HandleFunc("/user/signup", controllers.Register)
	http.HandleFunc(systemConf.Avatar.Uri, controllers.ImgShow)

	servConf := systemConf.Serv
	http.ListenAndServe(
		fmt.Sprintf("%s:%d", servConf.Host, servConf.Port),
		nil)
}
