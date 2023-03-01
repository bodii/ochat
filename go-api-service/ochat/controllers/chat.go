package controllers

import (
	"log"
	"ochat/service"

	"golang.org/x/net/websocket"
)

// chat the data received on the WebSocket.
func Chat(ws *websocket.Conn) {
	log.Printf("location: %v\n", ws.Config().Location)
	// log.Printf("origin: %v\n", ws.Config().Origin)
	r := ws.Request()
	// verify user legal
	user, code, errStr := service.NewUserServ().CheckUserRequestLegal(r)
	if errStr != "" {
		log.Println("websocket conn before get userinfo:", errStr)
		service.WsRespFailute(ws, code, errStr)
		return
	}
	log.Println(user)

	client := service.NewWsCline(ws, user.UserId)

	client.SendSystemMessage(user.UserId, "hello, welcome you")

	log.Printf("webSocket 与客户端建立连接: %#v  senderId: %d\n", client.Addr, user.UserId)

	client.Start()
}
