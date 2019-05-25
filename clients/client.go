package clients

import (
	"github.com/armatusmiles/gospeak-chat-engine/db/models"
	"github.com/gorilla/websocket"
)

type RegisteredChatClient struct {
	// The websocket connection.
	conn *websocket.Conn
}

func (c *RegisteredChatClient) SendMessage(msg *models.Message) {

}
