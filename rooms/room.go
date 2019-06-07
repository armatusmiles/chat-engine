package rooms

import (
	"github.com/gospeak/chat-engine/clients"
	"github.com/gospeak/chat-engine/rooms/broadcast"
)

type CommonChatRoom struct {
	broadcaster *broadcast.Broadcaster
	clients     []clients.ChatClient
}
