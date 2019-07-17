package clients

import (
	"sync"
)

// Threadsafe list
type ClientList struct {
	clients map[string]ChatClient
	mutex   *sync.RWMutex
}

func NewChatClientList() *ClientList {
	cl := &ClientList{
		clients: make(map[string]ChatClient),
		mutex:   &sync.RWMutex{},
	}
	return cl
}

// Don't forget. Count items in map can be changed after return value
func (cl *ClientList) Count() int {
	cl.mutex.RLock()
	len := len(cl.clients)
	cl.mutex.RUnlock()
	return len
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
	return true
}
