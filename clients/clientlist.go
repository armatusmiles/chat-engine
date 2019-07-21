package clients

import (
	"sync"
	"sync/atomic"
)

// Threadsafe list
type ClientList struct {
	counter uint32
	mutex   *sync.RWMutex
	clients map[string]ChatClient
}

func NewChatClientList() *ClientList {
	cl := &ClientList{
		clients: make(map[string]ChatClient),
		mutex:   &sync.RWMutex{},
		counter: 0,
	}
	return cl
}

// Don't forget. Count items in map can be changed after return value
func (cl *ClientList) Count() uint32 {
	return atomic.LoadUint32(&cl.counter)
}

// add only unique client
func (cl *ClientList) Add(client ChatClient) bool {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	_, exists := cl.clients[client.sessionID]
	if exists { // client already exists in map
		return false
	}

	cl.clients[client.sessionID] = client
	atomic.AddUint32(&cl.counter, 1)

	return true
}

func (cl *ClientList) RemoveBySessionID(ID string) bool {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	_, exists := cl.clients[ID]
	if !exists { // client not found in map
		return false
	}

	delete(cl.clients, ID)
	atomic.AddUint32(&cl.counter, ^uint32(0))

	return true
}

func (cl *ClientList) IsExists(ID string) bool {
	cl.mutex.RLock()
	defer cl.mutex.RUnlock()
	_, exists := cl.clients[ID]
	return exists
}
