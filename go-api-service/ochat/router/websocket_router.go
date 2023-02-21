package router

import (
	"net/http"
	"ochat/controllers"
	"ochat/service"

	"golang.org/x/net/websocket"
)

func WebsocketInit() {
	service.InitClientPool()

	http.Handle("/chat", websocket.Handler(controllers.Chat))
}
