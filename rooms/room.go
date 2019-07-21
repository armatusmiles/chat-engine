package rooms

import (
	"github.com/gorilla/websocket"
	"github.com/gospeak/chat-engine/clients"
)

type GeneralChatRoom struct {
	clients     *clients.ClientList
	broadcaster *clients.Broadcaster
}

func NewGeneralChatRoom() *GeneralChatRoom {
	gcr := &GeneralChatRoom{
		clients: clients.NewChatClientList(),
	}
	gcr.broadcaster = clients.NewBroadcaster(gcr.clients)
	return gcr
}

func (gcr *GeneralChatRoom) AddClient(conn *websocket.Conn, sessionID string) bool {
	client := clients.NewChatClient(conn, sessionID, gcr.broadcaster.ReadCh)
	return gcr.clients.Add(*client)
}

func (gcr *GeneralChatRoom) CountClients() uint32 {
	return gcr.clients.Count()
}

func (gcr *GeneralChatRoom) RemoveClientBySessionID(ID string) bool {
	return gcr.clients.RemoveBySessionID(ID)
}
