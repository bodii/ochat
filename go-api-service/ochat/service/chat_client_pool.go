package service

import (
	"sync"
)

// 连接池
var clientPool *PoolT

func InitClientPool() {
	clientPool = newPool()
}

type PoolT struct {
	clients map[int64]*ClientT
	lock    sync.RWMutex
}

func newPool() *PoolT {
	return &PoolT{
		clients: make(map[int64]*ClientT),
	}
}

// 设置用户Client到链接池
func setUserClient(userId int64, client *ClientT) bool {
	clientPool.lock.Lock()
	defer clientPool.lock.Unlock()

	clientPool.clients[userId] = client

	return true
}

func getUserClient(userId int64) (*ClientT, bool) {
	clientPool.lock.RLock()
	defer clientPool.lock.RUnlock()

	if client, ok := clientPool.clients[userId]; ok {
		return client, ok
	}

	return nil, false
}

func delUserClient(userId int64) (*ClientT, bool) {
	clientPool.lock.Lock()
	defer clientPool.lock.Unlock()

	if client, ok := clientPool.clients[userId]; ok {
		delete(clientPool.clients, userId)
		return client, true
	}

	return nil, false
}

func getclients() map[int64]*ClientT {
	clientPool.lock.RLock()
	defer clientPool.lock.RUnlock()

	return clientPool.clients
}

// getclientsUsersId
func getclientsUsersId() []int64 {
	clientPool.lock.RLock()
	defer clientPool.lock.RUnlock()

	userIds := make([]int64, len(clientPool.clients))
	i := 0
	for userid := range clientPool.clients {
		userIds[i] = userid
		i++
	}

	return userIds
}
