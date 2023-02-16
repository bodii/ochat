package service

import (
	"fmt"
	"ochat/comm"
	"ochat/models"
	"runtime/debug"
	"sync"
	"time"

	"golang.org/x/exp/slog"
	"golang.org/x/net/websocket"
)

type Client struct {
	SocketConn *websocket.Conn      // socket连接
	DataQueue  chan *models.Message // 待发送的数据
	Addr       string               // 客户端地址
	UserId     int64                // 用户id
	Type       string               // 客户端类型
	Wg         *sync.WaitGroup      // 上下文管理器
}

func NewWsCline(ws *websocket.Conn, userId int64, wg *sync.WaitGroup) *Client {
	wg.Add(1)

	client := &Client{
		SocketConn: ws,
		DataQueue:  make(chan *models.Message, 200),
		Addr:       ws.RemoteAddr().String(),
		UserId:     userId,
		Wg:         wg,
	}

	// 添加到连接池
	SetUserClient(userId, client)

	return client
}

// 接收协程
func (c *Client) Receive() {
	fmt.Println("start receive")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("receive stop", string(debug.Stack()), r)
		}
	}()
	defer func() {
		fmt.Println("读取客户端数据 关闭send", c)
		c.Close()
	}()

	for {
		// 接收数据
		var data models.Message
		// websocket.JSON.Receive(c.SocketConn, &data)
		err := websocket.JSON.Receive(c.SocketConn, &data)
		if err != nil {
			Log.Error("receive: ", err)
			return
		}

		if data.ReqId != "" {
			Log.Info("recv <<<< :",
				slog.String("client addr", c.Addr),
				slog.Any("receive data", data))

			// 发送数据
			c.dispatch(&data)
		}
	}
}

// 发送协程
func (c *Client) Send() {
	fmt.Println("start send")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("send stop", string(debug.Stack()), r)
		}
	}()

	defer func() {
		c.Close()
		fmt.Println("Client发送数据 defer", c)
	}()

	for {
		data := <-c.DataQueue
		if data != nil && data.ReqId != "" {
			Log.Info("send >>>> ",
				slog.String("client addr", c.Addr),
				slog.Any("receive data", data))

			err := websocket.JSON.Send(c.SocketConn, &comm.ResType{
				Code: 1,
				Msg:  "ok",
				Data: data,
			})

			if err != nil {
				Log.Error("send", err)
				websocket.JSON.Send(c.SocketConn, &comm.ResType{
					Code: 1001,
					Msg:  "send failure",
				})

				c.Close()
				return
			}

			fmt.Printf("\n\n")
		}
	}

}

// 处理接收到的数据
func (c *Client) dispatch(data *models.Message) {
	senderClient := GetUserClient(data.ReceiverId)
	fmt.Printf("receiver: %#v\n", senderClient)
	// 根据message的模式处理
	switch data.Mode {
	case models.MESSAGE_MODE_SINGLE:
		c.SendMessage(data)
	case models.MESSAGE_MODE_GROUP:
		// 是否还存在连接
		c.SendGroupMessage(data)
	}
}

// 发送信息
func (c *Client) SendMessage(data *models.Message) {
	fmt.Println("in sendMessage: ", data)
	// 发送信息自已要收到一条
	go func() {
		c.DataQueue <- data
	}()

	// 如果不是系统发送的
	if data.SenderId != data.ReceiverId && data.SenderId > 0 {
		// 接收者收到一条
		receiverCilent := GetUserClient(data.ReceiverId)
		if receiverCilent != nil {
			go func() {
				receiverCilent.DataQueue <- data
			}()
		}
	}

	fmt.Println("[[ send end ]]")
}

// 发送群信息
func (c *Client) SendGroupMessage(data *models.Message) {
	for _, client := range GetClients() {
		client.SendMessage(data)
	}
}

// 发送系统信息
func (c *Client) SendSystemMessage(userId int64, msg string) {
	message := &models.Message{
		ReqId:           "2001",
		SenderId:        0,
		ReceiverId:      userId,
		Mode:            models.MESSAGE_MODE_SINGLE,
		Type:            models.MESSAGE_TYPE_SYSTEM,
		Content:         msg,
		SenderStatus:    1,
		ReceiverStatus:  1,
		CreatedAt:       time.Now(),
		SenderUpdatedAt: time.Now(),
	}

	fmt.Println("in sendSystem: exec sendMessage")
	c.SendMessage(message)
}

func WsRespFailute(ws *websocket.Conn, code int, msg string) {
	websocket.JSON.Send(ws, &comm.ResType{
		Code: code,
		Msg:  msg,
	})
	ws.Close()
}

func (c *Client) Close() {
	c.SocketConn.Close()

	c.Wg.Done()
	// close(c.DataQueue)

	DelUserClient(c.UserId)
}
