package router

import (
	"net/http"
	"ochat/controllers"

	"golang.org/x/net/websocket"
)

func WebsocketInit() {
	http.Handle("/chat", websocket.Handler(controllers.Chat))
}
