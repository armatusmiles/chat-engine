package clients

import (
	dbm "github.com/gospeak/protorepo/dbmanage"
)

type Broadcaster struct {
	clientList *ClientList
	ReadCh     chan dbm.ChatMessage
	cancelCh   chan struct{}
}

func NewBroadcaster(clients *ClientList) *Broadcaster {
	rb := &Broadcaster{
		clientList: clients,
		ReadCh:     make(chan dbm.ChatMessage),
		cancelCh:   make(chan struct{}),
	}
	go rb.chReader()
	return rb
}

func (rb *Broadcaster) chReader() {
	for {
		select {
		case <-rb.cancelCh:
			return
		case msg := <-rb.ReadCh:
			rb.clientList.mutex.RLock()
			for _, el := range rb.clientList.clients {
				el.SendMessage(&msg)
			}
			rb.clientList.mutex.RUnlock()
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
