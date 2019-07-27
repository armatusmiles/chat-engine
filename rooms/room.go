package rooms

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/gospeak/chat-engine/clients"
	bc "github.com/gospeak/chat-engine/rooms/broadcaster"
)

type GeneralChatRoom struct {
	broadcaster *bc.Broadcaster
	mutex       *sync.RWMutex
	clients     *clients.ClientList
}

func NewGeneralChatRoom() *GeneralChatRoom {
	gcr := &GeneralChatRoom{
		clients: clients.NewChatClientList(),
		mutex:   &sync.RWMutex{},
	}
	gcr.broadcaster = bc.NewBroadcaster(gcr.clients, gcr.mutex)
	return gcr
}

// You should think about close websocket.Conn if this function return false
func (gcr *GeneralChatRoom) AddClient(conn *websocket.Conn, sessionID string) bool {
	gcr.mutex.Lock()
	defer gcr.mutex.Unlock()
	if gcr.clients.IsExists(sessionID) {
		return false
	}
	client := clients.NewChatClient(conn, sessionID, gcr.broadcaster.ReadCh)
	return gcr.clients.Add(*client)
}

// Attention! Quantity items in map can be changed after return value
// Race condition is possible
func (gcr *GeneralChatRoom) CountClients() uint32 {
	return gcr.clients.Count()
}

func (gcr *GeneralChatRoom) RemoveClientBySessionID(ID string) bool {
	gcr.mutex.Lock()
	defer gcr.mutex.Unlock()
	return gcr.clients.RemoveBySessionID(ID)
}
