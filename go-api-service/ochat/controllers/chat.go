package controllers

import (
	"fmt"
	"log"
	"ochat/comm"
	"ochat/service"
	"sync"

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
		WsRespFailute(ws, code, errStr)
		return
	}
	fmt.Println(user)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	client := service.NewWsCline(ws, user.UserId, wg)

	client.SendSystemMessage(user.UserId, "hello, welcome you")

	fmt.Printf("\nwebSocket 与客户端建立连接: %#v  senderId: %d\n\n", client.Addr, user.UserId)

	// 接收
	go client.Receive()
	// 发送
	go client.Send()

	wg.Wait()

	fmt.Println("\n\n [[[   end   ]]]")
}

func WsRespFailute(ws *websocket.Conn, code int, msg string) {
	websocket.JSON.Send(ws, &comm.ResType{
		Code: code,
		Msg:  msg,
	})
	ws.Close()
}
