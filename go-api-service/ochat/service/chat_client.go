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

type ClientT struct {
	SocketConn *websocket.Conn      // socket连接
	DataQueue  chan *models.Message // 待发送的数据
	Addr       string               // 客户端地址
	UserId     int64                // 用户id
	Type       string               // 客户端类型
	Wg         *sync.WaitGroup      // 上下文管理器
}

// 接收信息
type ReceiveMessageT struct {
	From    int64  `json:"from" form:"from"`
	To      int64  `json:"to,omitempty" form:"to"`
	Mode    int    `json:"mode" form:"mode"`
	Type    int    `json:"type" form:"type"`
	Content string `json:"content,omitempty" form:"content"`
	Pic     string `json:"pic,omitempty" form:"pic"`
	Url     string `json:"url,omitempty" form:"url"`
	About   string `json:"about,omitempty" form:"about"`
	Amount  int    `json:"amount,omitempty" form:"amount"`
}

func NewWsCline(ws *websocket.Conn, userId int64) *ClientT {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	client := &ClientT{
		SocketConn: ws,
		DataQueue:  make(chan *models.Message, 800),
		Addr:       ws.RemoteAddr().String(),
		UserId:     userId,
		Wg:         wg,
	}

	// 添加到连接池
	setUserClient(userId, client)

	return client
}

// 接收协程
func (c *ClientT) Receive() {
	fmt.Println("start receive")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("receive stop", string(debug.Stack()), r)
		}
	}()

	defer c.Close()

	for {
		// 接收数据
		var data ReceiveMessageT
		err := websocket.JSON.Receive(c.SocketConn, &data)
		fmt.Println("data===========:", data)
		if err != nil {
			slog.Error("receive: ", err, slog.Any("receive data:", data))
			return
		}

		if data.Content == "" {
			continue
		}

		slog.Info("recv <<<< :",
			slog.String("client addr", c.Addr),
			slog.Any("receive data", data))

		// 将接收到的数据转换成models.Message
		message := receiveToMessage(data)

		// 发送数据
		c.dispatch(message)

	}
}

// 发送协程
func (c *ClientT) Send() {
	fmt.Println("start send")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("send stop", string(debug.Stack()), r)
		}
	}()

	defer c.Close()

	for {
		data := <-c.DataQueue
		fmt.Println("+++data+++++++", data)

		slog.Info("send >>>> ",
			slog.String("client addr", c.Addr),
			slog.Any("receive data", data))

		err := websocket.JSON.Send(c.SocketConn, &comm.ResType{
			Code: 1,
			Msg:  "ok",
			Data: data,
		})

		fmt.Println()

		if err == nil {
			continue
		}

		slog.Error("send", err)
		WsRespFailute(c.SocketConn, 1001, "send failure")

		// 如果有错误，结束联系
		c.Close()
		return
	}
}

// 处理接收到的数据
func (c *ClientT) dispatch(data *models.Message) {
	// 获取接收者Client
	senderClient, ok := getUserClient(data.ReceiverId)
	if !ok {
		fmt.Println("\n--||||--receiver not exists")
	}

	fmt.Printf("\n--||||--receiver: %#v\n", senderClient)

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
func (c *ClientT) SendMessage(data *models.Message) {
	fmt.Println("in sendMessage: ", data)
	// 不是系统发送，发送信息自已要收到一条
	if data.SenderId != 0 {
		go func() {
			c.DataQueue <- data
		}()
	}

	// 接收者收到一条
	receiverCilent, ok := getUserClient(data.ReceiverId)
	fmt.Println("\n {{{ receiverClient }}}", receiverCilent)
	if ok && receiverCilent != nil {
		go func() {
			receiverCilent.DataQueue <- data
		}()
	}

	// 保存数据到数据库
	go func() {
		data, err := NewMessageServ().Save(data)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println("\n saved =====,===", data)
	}()

	fmt.Println("[[ send end ]]")
}

// 发送群信息
func (c *ClientT) SendGroupMessage(data *models.Message) {
	groups, err := NewGroupContactServ().GetMembers(data.GroupId)
	if err != nil {
		WsRespFailute(c.SocketConn, 1001, "send group message failure, group info not exists")
		return
	}

	for _, member := range groups {
		data.ReceiverId = member.UserId
		c.SendMessage(data)
	}
}

// 发送系统信息
func (c *ClientT) SendGroupSystemMessage(GroupId int64, msg string) {
	message := &models.Message{
		SenderId:       0,
		GroupId:        GroupId,
		Mode:           models.MESSAGE_MODE_GROUP,
		Type:           models.MESSAGE_TYPE_SYSTEM,
		Content:        msg,
		SenderStatus:   1,
		ReceiverStatus: 1,
	}

	fmt.Println("in sendGroupSystem: exec sendGroupMessage")
	c.SendGroupMessage(message)
}

// 发送系统信息
func (c *ClientT) SendSystemMessage(userId int64, msg string) {
	message := &models.Message{
		SenderId:       0,
		ReceiverId:     userId,
		Mode:           models.MESSAGE_MODE_SINGLE,
		Type:           models.MESSAGE_TYPE_SYSTEM,
		Content:        msg,
		SenderStatus:   1,
		ReceiverStatus: 1,
	}

	fmt.Println("in sendSystem: exec sendMessage")
	c.SendMessage(message)
}

// 开启连接
func (c *ClientT) Start() {
	// 接收消息通道
	go c.Receive()

	// 发送消息通道
	go c.Send()

	// 线程等待
	c.Wg.Wait()
}

// 关闭连接
func (c *ClientT) Close() {
	c.SocketConn.Close()

	// 结束协程
	c.Wg.Done()

	delUserClient(c.UserId)

	fmt.Println("\n\n [[[   end   ]]]")
}

// 错误时的返回信息
func WsRespFailute(ws *websocket.Conn, code int, msg string) {
	websocket.JSON.Send(ws, &comm.ResType{
		Code: code,
		Msg:  msg,
	})

	ws.Close()
}

func receiveToMessage(data ReceiveMessageT) *models.Message {
	now := time.Now()

	message := &models.Message{
		SenderId:          data.From,
		Mode:              data.Mode,
		Type:              data.Type,
		Content:           data.Content,
		Pic:               data.Pic,
		Url:               data.Url,
		About:             data.About,
		Amount:            data.Amount,
		SenderStatus:      1,
		ReceiverStatus:    1,
		CreatedAt:         now,
		SenderUpdatedAt:   now,
		ReceiverUpdatedAt: now,
	}

	if data.Mode == models.MESSAGE_MODE_SINGLE {
		message.ReceiverId = data.To
	} else if data.Mode == models.MESSAGE_MODE_GROUP {
		message.GroupId = data.To
	}

	return message
}
