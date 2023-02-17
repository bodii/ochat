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
	ClientPool = newPool()
}

func InitLog() {
	textHandler := slog.NewTextHandler(os.Stdout)
	Log = slog.New(textHandler)
}

func newPool() *Pool {
	return &Pool{
		Clients: make(map[int64]*Client),
		Lock:    &sync.RWMutex{},
	}
}

// 设置用户Client到链接池
func setUserClient(userId int64, client *Client) bool {
	ClientPool.Lock.Lock()
	defer ClientPool.Lock.Unlock()

	ClientPool.Clients[userId] = client

	return true
}

func getUserClient(userId int64) (*Client, bool) {
	ClientPool.Lock.RLock()
	defer ClientPool.Lock.RUnlock()

	if client, ok := ClientPool.Clients[userId]; ok {
		return client, ok
	}

	return nil, false
}

func delUserClient(userId int64) (*Client, bool) {
	ClientPool.Lock.Lock()
	defer ClientPool.Lock.Unlock()

	if client, ok := ClientPool.Clients[userId]; ok {
		delete(ClientPool.Clients, userId)
		return client, true
	}

	return nil, false
}

func getClients() map[int64]*Client {
	ClientPool.Lock.RLock()
	defer ClientPool.Lock.RUnlock()

	return ClientPool.Clients
}

// getClientsUsersId
func getClientsUsersId() []int64 {
	ClientPool.Lock.RLock()
	defer ClientPool.Lock.RUnlock()

	userIds := make([]int64, 0)
	for userid := range ClientPool.Clients {
		userIds = append(userIds, userid)
	}

	return userIds
}
