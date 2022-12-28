package controllers

import (
	"fmt"
	"log"
	"ochat/comm"
	"ochat/models"
	"ochat/service"
	"strconv"
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
	"golang.org/x/net/websocket"
)

var rwMut sync.RWMutex

type Node struct {
	WsConn *websocket.Conn
	// 并行转串行，
	DataQueue chan models.Message
	GroupSets mapset.Set[int64]
}

var clientMap map[int64]*Node = make(map[int64]*Node)

// chat the data received on the WebSocket.
func Chat(ws *websocket.Conn) {
	// log.Printf("location: %v\n", ws.Config().Location)
	// log.Printf("origin: %v\n", ws.Config().Origin)
	r := ws.Request()
	userIdStr := r.FormValue("user_id")
	token := r.FormValue("token")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		websocket.JSON.Send(ws, &comm.R{
			Code: 1001,
			Msg:  "user info err",
		})
		ws.Close()
		return
	}
	if token == "" {
		websocket.JSON.Send(ws, &comm.R{
			Code: 1001,
			Msg:  "token info err",
		})
		ws.Close()
		return
	}

	isValida := service.NewUserServ().CheckToken(userId, token)
	if !isValida {
		websocket.JSON.Send(ws, &comm.R{
			Code: 1001,
			Msg:  "token valida failure",
		})
		ws.Close()
		return
	}

	node := &Node{
		WsConn:    ws,
		DataQueue: make(chan models.Message),
		GroupSets: mapset.NewSet[int64](),
	}

	rwMut.Lock()
	clientMap[userId] = node
	rwMut.Unlock()

	// 接收
	go recvproc(node)
}

// 处理
func dispatch(data models.Message) (err error) {
	// 根据message的模式处理
	switch data.Mode {
	case models.MESSAGE_MODE_SINGLE:
		sendMessage(data.ReceiverId, data)
	case models.MESSAGE_MODE_GROUP:
		// 是否还存在连接
		sendGroupMessage(data.ReceiverId, data)
	}

	return nil
}

var rowlocker sync.RWMutex

// 发送信息
func sendMessage(userId int64, data models.Message) {
	rowlocker.RLock()
	node, ok := clientMap[userId]
	rowlocker.RUnlock()
	if ok {
		node.DataQueue <- data
	}
}

// 发送群信息
func sendGroupMessage(userId int64, data models.Message) {
	for _, v := range clientMap {
		if v.GroupSets.Contains(userId) {
			v.DataQueue <- data
		}
	}
}

// 接收协程
func recvproc(node *Node) {
	for {
		var data models.Message
		if err := websocket.JSON.Receive(node.WsConn, &data); err != nil {
			fmt.Println(err.Error())
			return
		}

		go dispatch(data)

		// 对data进行进一步的处理
		fmt.Printf("recv<-:%v \n", data)
	}
}

// 发送协程
func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := websocket.JSON.Send(node.WsConn, data)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}
