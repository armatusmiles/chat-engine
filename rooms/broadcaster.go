package rooms

import (
	"github.com/gospeak/chat-engine/clients"
	"github.com/gospeak/db-manage-service/redis/models"
	"github.com/gospeak/protorepo/dbmanage"
)

type IRoomBroadcaster interface {
	BroadcastMessage(msg *models.Message)
}

type RoomBroadcaster struct {
	clients []clients.ChatClient
	ReadCh  chan []dbmanage.ChatMessage
}

func (rb *RoomBroadcaster) BroadcastMessage(msg *dbmanage.ChatMessage) {
	// TODO send message to all client
}
