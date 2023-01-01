package controllers

import (
	"fmt"
	"ochat/comm"
	"ochat/models"
	"ochat/service"
	"strconv"
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
	"golang.org/x/net/websocket"
)

var rwMut sync.RWMutex

type Client struct {
	WsConn *websocket.Conn
	// 并行转串行，
	DataQueue chan *models.Message
	GroupSets mapset.Set[int64]
}

var clientPool map[int64]*Client = make(map[int64]*Client)

// chat the data received on the WebSocket.
func Chat(ws *websocket.Conn) {
	// log.Printf("location: %v\n", ws.Config().Location)
	// log.Printf("origin: %v\n", ws.Config().Origin)
	r := ws.Request()
	userIdStr := r.FormValue("user_id")
	token := r.FormValue("token")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		WsRespFailute(ws, 1001, "user info err")
		return
	}
	if token == "" {
		WsRespFailute(ws, 1001, "token info err")
		return
	}

	isValida := service.NewUserServ().CheckToken(userId, token)
	if !isValida {
		WsRespFailute(ws, 1001, "token valida failure")
		return
	}

	client := NewWsCline(ws, userId)

	// 接收
	go client.recvProc()
	// 发送
	go client.sendProc()

	client.sendSystemMessage(userId, "hello, welcome you")

	fmt.Printf("\nwebSocket 与客户端建立连接: %#v  senderId: %d\n\n", ws.Request().RemoteAddr, userId)
	client.sendProc()
}

func NewWsCline(ws *websocket.Conn, userId int64) *Client {
	client := &Client{
		WsConn:    ws,
		DataQueue: make(chan *models.Message),
		GroupSets: mapset.NewSet[int64](),
	}

	SetUserClient(userId, client)

	return client
}

// 设置用户Client到链接池
func SetUserClient(userId int64, client *Client) {
	rwMut.Lock()
	clientPool[userId] = client
	rwMut.Unlock()
}

func GetUserClient(userId int64) *Client {
	rwMut.RLock()
	client := clientPool[userId]
	rwMut.RUnlock()

	return client
}

func WsRespFailute(ws *websocket.Conn, code int, msg string) {
	websocket.JSON.Send(ws, &comm.R{
		Code: code,
		Msg:  msg,
	})
	ws.Close()
}

// 接收协程
func (c *Client) recvProc() {
	for {
		// 接收数据
		var data models.Message
		err := websocket.JSON.Receive(c.WsConn, &data)
		if err != nil {
			fmt.Println("recvProc: ", err.Error())
			return
		}

		fmt.Printf("%s recv<-: %#v\n", c.WsConn.Request().RemoteAddr, data)

		// 发送数据
		c.dispatch(data)
	}
}

// 发送协程
func (c *Client) sendProc() {
	for {
		data := <-c.DataQueue
		fmt.Printf("send-> %s: %#v\n", c.WsConn.Request().RemoteAddr, data)
		err := websocket.JSON.Send(c.WsConn, &comm.R{
			Code: 1,
			Msg:  "ok",
			Data: data,
		})
		if err != nil {
			fmt.Printf("sendProc: %#v\n", err.Error())
			websocket.JSON.Send(c.WsConn, &comm.R{
				Code: 1001,
				Msg:  "send failure",
			})

			c.Close()
			return
		}
	}
}

// 处理
func (c *Client) dispatch(data models.Message) {
	senderClient := GetUserClient(data.ReceiverId)
	fmt.Printf("receiver: %#v\n", senderClient)
	// 根据message的模式处理
	switch data.Mode {
	case models.MESSAGE_MODE_SINGLE:
		c.sendMessage(data)
		// sender.sendMessage(data)
	case models.MESSAGE_MODE_GROUP:
		// 是否还存在连接
		c.sendGroupMessage(data)
		// sender.sendGroupMessage(data)
	}
}

// 发送信息
func (c *Client) sendMessage(data models.Message) {
	// 发送信息自已要收到一条
	c.DataQueue <- &data
	// 如果不是系统发送的
	if data.SenderId != data.ReceiverId {
		// 接收者收到一条
		receiverCilent := GetUserClient(data.ReceiverId)
		if receiverCilent != nil {
			receiverCilent.DataQueue <- &data
		}
	}
}

// 发送群信息
func (c *Client) sendGroupMessage(data models.Message) {
	for _, v := range c.GroupSets.ToSlice() {
		fmt.Println(v)
	}
}

// 发送系统信息
func (c *Client) sendSystemMessage(userId int64, msg string) {
	c.sendMessage(models.Message{
		SenderId:   userId,
		ReceiverId: userId,
		Mode:       1,
		Type:       10,
		Content:    msg,
	})
}

func (c *Client) Close() {
	c.WsConn.Close()
	close(c.DataQueue)
}
