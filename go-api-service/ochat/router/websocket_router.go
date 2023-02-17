package router

import (
	"net/http"
	"ochat/controllers"
	"ochat/service"

	"golang.org/x/net/websocket"
)

func WebsocketInit() {
	service.InitClientPool()
	service.InitLog()

	http.Handle("/chat", websocket.Handler(controllers.Chat))
}
