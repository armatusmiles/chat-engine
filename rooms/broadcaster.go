package main

import (
	"github.com/armatusmiles/gospeak-chat-engine/db/models"
)

type IRoomBroadcaster interface {
	BroadcastMessage(msg *models.Message)
}

type RoomBroadcaster struct {
	clients *[]ChatClient
}

func (rb *RoomBroadcaster) BroadcastMessage(msg *models.Message) {
	// TODO send message to all client
}
