package broadcaster

import (
	"sync"

	"github.com/gospeak/chat-engine/clients"
	dbm "github.com/gospeak/protorepo/dbmanage"
)

type Broadcaster struct {
	ReadCh      chan dbm.ChatMessage
	cancelCh    chan struct{}
	clientMutex *sync.RWMutex
	clientList  *clients.ClientList
}

func NewBroadcaster(clients *clients.ClientList, cm *sync.RWMutex) *Broadcaster {
	rb := &Broadcaster{
		clientList:  clients,
		clientMutex: cm,
		ReadCh:      make(chan dbm.ChatMessage),
		cancelCh:    make(chan struct{}),
	}
	go rb.chReader()
	return rb
}

// Reads client messages from the channel.
// Broadcasts message to all members in the list.
func (rb *Broadcaster) chReader() {
	for {
		select {
		case <-rb.cancelCh:
			return
		case msg := <-rb.ReadCh:
			rb.clientMutex.RLock()
			rb.clientList.SendMessageToAll(&msg)
			rb.clientMutex.RUnlock()
		}
	}
}

func (rb *Broadcaster) BroadcastMessage(msg *dbm.ChatMessage) {
	rb.ReadCh <- *msg
}

// shutdown chReader gorutione
func (rb *Broadcaster) Close() {
	rb.cancelCh <- struct{}{}
}
