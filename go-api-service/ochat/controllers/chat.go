package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ochat/bootstrap"
	"ochat/comm"
	"ochat/models"
	"ochat/service"
	"strconv"
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gorilla/websocket"
)

var rwMut sync.RWMutex

type Node struct {
	Conn *websocket.Conn
	// 并行转串行，
	DataQueue chan []byte
	GroupSets mapset.Set[int64]
}

var clientMap map[int64]*Node = make(map[int64]*Node)

func Chat(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userIdStr := r.FormValue("user_id")
	token := r.FormValue("token")

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		comm.Res(w, 1001, err.Error(), nil)
		return
	}

	if token == "" {
		comm.Res(w, 1002, "token is empty", nil)
		return
	}

	isValida := checkToken(userId, token)

	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isValida
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket err:%v", err)
		return
	}

	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: mapset.NewSet[int64](),
	}

	rwMut.Lock()
	clientMap[userId] = node
	rwMut.Unlock()

	// 发送
	go sendProc(node)

	sendMessage(userId, []byte("Hello, world"))

}

// 处理
func dispatch(data []byte) (err error) {
	// 解析data为message
	message := models.Message{}

	err = json.Unmarshal(data, &message)
	if err != nil {
		return err
	}
	// 根据message的模式处理
	switch message.Mode {
	case models.MESSAGE_MODE_SINGLE:
		sendMessage(message.ReceiverId, data)
	case models.MESSAGE_MODE_GROUP:
		// 是否还存在连接
		sendGroupMessage(message.ReceiverId, data)
	}

	return nil
}

var rowlocker sync.RWMutex

// 发送信息
func sendMessage(userId int64, msg []byte) {
	rowlocker.RLock()
	node, ok := clientMap[userId]
	rowlocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}

// 发送群信息
func sendGroupMessage(userId int64, msg []byte) {
	for _, v := range clientMap {
		if v.GroupSets.Contains(userId) {
			v.DataQueue <- msg
		}
	}
}

func checkToken(userId int64, token string) bool {
	userServ := &service.UserService{
		DB: bootstrap.DB_Engine,
	}

	return userServ.CheckToken(userId, token)
}

// 接收协程
func recvproc(node *Node) {
	for {
		msgType, data, err := node.Conn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			return
		}

		go dispatch(data)

		// 对data进行进一步的处理
		fmt.Printf("recv<-:%v type:\n", data, msgType)
	}
}

// 发送协程
func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}
