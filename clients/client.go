package clients

import (
	"github.com/sirupsen/logrus"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	dbm "github.com/gospeak/protorepo/dbmanage"
)

var log = logrus.New()

type ChatClient struct {
	// The websocket connection.
	conn       *websocket.Conn
	session_id string

	readCh chan<- dbm.ChatMessage
}

func NewChatClient(conn *websocket.Conn, session_id string,
	readCh chan<- dbm.ChatMessage) *ChatClient {
	cc := &ChatClient{conn, session_id, readCh}
	go cc.readThread()
	return cc
}

// readThread reads messages from socket, unmarhals to ChatMessage and send to readCh
func (c *ChatClient) readThread() {
	for {
		mt, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Error(err)
			return
		}
		if mt != websocket.BinaryMessage {
			log.Error("Must be binary (proto) message!")
		}
		var msg dbm.ChatMessage
		err = proto.Unmarshal(message, &msg)
		if err != nil {
			log.Error("Fail unpack message")
		}
		c.readCh <- msg
	}
}

func (c *ChatClient) CloseConnection() {
	c.conn.Close()
}

// SendMessage sends message to client
func (c *ChatClient) SendMessage(msg *dbm.ChatMessage) {
	// todo write in channel
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Error("Marshaling error: ", err)
		return
	}
	c.conn.WriteMessage(websocket.BinaryMessage, data)
}
