package clients

import (
	"github.com/sirupsen/logrus"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/gospeak/protorepo/dbmanage"
)

var log = logrus.New()

type ChatClient struct {
	// The websocket connection.
	conn       *websocket.Conn
	session_id string
}

func NewChatClient(conn *websocket.Conn, session_id string) *ChatClient {
	return &ChatClient{conn, session_id}
}

// SendMessage sends message to client
func (c *ChatClient) SendMessage(msg *dbmanage.ChatMessage) {
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Error("Marshaling error: ", err)
		return
	}
	c.conn.WriteMessage(websocket.BinaryMessage, data)
}
