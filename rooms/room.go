package rooms

import "github.com/armatusmiles/gospeak-chat-engine/clients"

type FreeChatRoom struct {
	broadcaster *RoomBroadcaster
	clients     *[]clients.ChatClient
}
