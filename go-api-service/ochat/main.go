package main

import (
	"fmt"
	"net/http"
	"ochat/bootstrap"
	"ochat/controllers"

	"golang.org/x/net/websocket"
)

// main func
func main() {
	bootstrap.Init()

	systemConf := bootstrap.SystemConf

	http.HandleFunc("/user/login", controllers.Login)
	http.HandleFunc("/user/signup", controllers.Register)
	http.HandleFunc("/user/avatar", controllers.ShowAvatar)
	http.HandleFunc("/user/avatar/upload", controllers.UpPicture)
	http.Handle("/chat", websocket.Handler(controllers.Chat))

	servConf := systemConf.Serv
	http.ListenAndServe(
		fmt.Sprintf("%s:%d", servConf.Host, servConf.Port),
		nil)
}
