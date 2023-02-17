package service

import (
	"os"
	"sync"

	"golang.org/x/exp/slog"
)

// 连接池
var ClientPool *Pool
var Log *slog.Logger

type Pool struct {
	Clients map[int64]*Client
	Lock    *sync.RWMutex
}

func InitClientPool() {
	ClientPool = NewPool()
}

func InitLog() {
	textHandler := slog.NewTextHandler(os.Stdout)
	Log = slog.New(textHandler)
}

func NewPool() *Pool {
	return &Pool{
		Clients: make(map[int64]*Client),
		Lock:    &sync.RWMutex{},
	}
}

// 设置用户Client到链接池
func SetUserClient(userId int64, client *Client) {
	ClientPool.Lock.Lock()
	defer ClientPool.Lock.Unlock()

	ClientPool.Clients[userId] = client
}

func GetUserClient(userId int64) *Client {
	ClientPool.Lock.RLock()
	defer ClientPool.Lock.RUnlock()

	client := ClientPool.Clients[userId]
	return client
}

func DelUserClient(userId int64) *Client {
	ClientPool.Lock.Lock()
	defer ClientPool.Lock.Unlock()

	if client, ok := ClientPool.Clients[userId]; ok {
		delete(ClientPool.Clients, userId)
		return client
	}

	return nil
}

func GetClients() map[int64]*Client {
	ClientPool.Lock.RLock()
	defer ClientPool.Lock.RUnlock()

	return ClientPool.Clients
}

func GetClientsUsersId() []int64 {
	ClientPool.Lock.RLock()
	defer ClientPool.Lock.RUnlock()

	userIds := make([]int64, 0)
	for userid := range ClientPool.Clients {
		userIds = append(userIds, userid)
	}

	return userIds
}
