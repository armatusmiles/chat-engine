package broadcast

import (
	"sync"

	"github.com/gospeak/chat-engine/clients"
	dbm "github.com/gospeak/protorepo/dbmanage"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Broadcaster struct {
	ReadCh chan dbm.ChatMessage

	clients      *[]clients.ChatClient
	clientsMutex *sync.RWMutex
	cancelCh     chan struct{}
}

func NewBroadcaster(clients *[]clients.ChatClient) *Broadcaster {
	rb := &Broadcaster{
		ReadCh:       make(chan dbm.ChatMessage),
		cancelCh:     make(chan struct{}),
		clientsMutex: &sync.RWMutex{},
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
			rb.clientsMutex.RLock()
			for _, el := range *rb.clients {
				el.SendMessage(&msg)
			}
			rb.clientsMutex.RUnlock()
		}
	}
}

func (rb *Broadcaster) UpdateListClients(clients *[]clients.ChatClient) {
	rb.clientsMutex.Lock()
	rb.clients = clients
	rb.clientsMutex.Unlock()
}

func (rb *Broadcaster) BroadcastMessage(msg *dbm.ChatMessage) {
	rb.ReadCh <- *msg
}

// shutdown chReader gorutione
func (rb *Broadcaster) Close() {
	rb.cancelCh <- struct{}{}
}
