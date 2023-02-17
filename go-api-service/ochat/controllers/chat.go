package controllers

import (
	"fmt"
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
		fmt.Println("websocket conn before get userinfo:", errStr)
		service.WsRespFailute(ws, code, errStr)
		return
	}
	fmt.Println(user)

	client := service.NewWsCline(ws, user.UserId)

	client.SendSystemMessage(user.UserId, "hello, welcome you")

	fmt.Printf("\nwebSocket 与客户端建立连接: %#v  senderId: %d\n\n", client.Addr, user.UserId)

	client.Start()
}
