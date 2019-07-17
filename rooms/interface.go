package rooms

import (
	"github.com/gorilla/websocket"
)

type ChatRoom interface {
	CountClients() uint32
	AddClient(conn *websocket.Conn, sessionID string) bool
	RemoveClientBySessionID(ID string) bool
}
