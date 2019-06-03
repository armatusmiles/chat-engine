package rooms

import "github.com/armatusmiles/gospeak-chat-engine/clients"

type IChatRoom interface {
	CountClients() uint32
	AddClient(*clients.ChatClient) bool
	RemoveClientById(ID uint64) bool
}
