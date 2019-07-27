package clients

import (
	"sync/atomic"

	dbm "github.com/gospeak/protorepo/dbmanage"
)

// Threadsafe list
type ClientList struct {
	counter uint32
	clients map[string]ChatClient
}

func NewChatClientList() *ClientList {
	cl := &ClientList{
		clients: make(map[string]ChatClient),
		counter: 0,
	}
	return cl
}

// Sends message to all clients from server
func (cl *ClientList) SendMessageToAll(msg *dbm.ChatMessage) {
	for _, el := range cl.clients {
		el.SendMessage(msg)
	}
}

// Attention! The number of elements in the map can be changed after returning the value
func (cl *ClientList) Count() uint32 {
	return atomic.LoadUint32(&cl.counter)
}

// add only unique client
func (cl *ClientList) Add(client ChatClient) bool {
	_, exists := cl.clients[client.sessionID]
	if exists { // client already exists in map
		return false
	}

	cl.clients[client.sessionID] = client
	atomic.AddUint32(&cl.counter, 1)

	return true
}

func (cl *ClientList) RemoveBySessionID(ID string) bool {
	_, exists := cl.clients[ID]
	if !exists { // client not found in map
		return false
	}

	delete(cl.clients, ID)
	atomic.AddUint32(&cl.counter, ^uint32(0))

	return true
}

func (cl *ClientList) IsExists(ID string) bool {
	_, exists := cl.clients[ID]
	return exists
}
