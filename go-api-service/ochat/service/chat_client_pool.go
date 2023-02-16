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
	Pool map[int64]*Client
	Lock sync.RWMutex
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
		Pool: make(map[int64]*Client),
		Lock: sync.RWMutex{},
	}
}

// 设置用户Client到链接池
func SetUserClient(userId int64, client *Client) {
	ClientPool.Lock.Lock()
	defer ClientPool.Lock.Unlock()

	ClientPool.Pool[userId] = client
}

func GetUserClient(userId int64) *Client {
	ClientPool.Lock.RLock()
	defer ClientPool.Lock.RUnlock()

	client := ClientPool.Pool[userId]
	return client
}

func DelUserClient(userId int64) *Client {
	ClientPool.Lock.Lock()
	defer ClientPool.Lock.Unlock()

	if client, ok := ClientPool.Pool[userId]; ok {
		delete(ClientPool.Pool, userId)
		return client
	}

	return nil
}

func GetClients() map[int64]*Client {
	ClientPool.Lock.RLock()
	defer ClientPool.Lock.RUnlock()

	return ClientPool.Pool
}

func GetClientsUsersId() []int64 {
	ClientPool.Lock.RLock()
	defer ClientPool.Lock.RUnlock()

	userIds := make([]int64, 0)
	for userid := range ClientPool.Pool {
		userIds = append(userIds, userid)
	}

	return userIds
}
