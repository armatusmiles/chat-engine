package rooms

import (
	"github.com/gospeak/chat-engine/clients"
	dbm "github.com/gospeak/protorepo/dbmanage"
)

type ChatRoom interface {
	CountClients() uint32
	AddClient(*clients.ChatClient) bool
	RemoveClientById(ID uint64) bool
}

type Broadcaster interface {
	BroadcastMessage(msg *dbm.ChatMessage)
}
