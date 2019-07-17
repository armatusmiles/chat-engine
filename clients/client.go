package clients

import (
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	dbm "github.com/gospeak/protorepo/dbmanage"
	log "github.com/sirupsen/logrus"
)

type ChatClient struct {
	// The websocket connection.
	conn      *websocket.Conn
	sessionID string

	readCh chan<- dbm.ChatMessage // broadcaster channel
}

func NewChatClient(conn *websocket.Conn, sessionID string,
	readCh chan<- dbm.ChatMessage) *ChatClient {
	cc := &ChatClient{conn, sessionID, readCh}
	go cc.readThread()
	return cc
}

func (c *ChatClient) GetSessionID() string {
	return c.sessionID
}

// readThread reads messages from socket, unmarhals to ChatMessage
// and send to broadcaster (readCh)
func (c *ChatClient) readThread() {
	defer c.conn.Close()
	for {
		mt, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Error(err)
			return
		}
		log.Debugf("Read message readThread() %s", message)
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

// SendMessage sends message from server to client
func (c *ChatClient) SendMessage(msg *dbm.ChatMessage) error {
	// todo write in channel
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Error("Marshaling error: ", err)
		return err
	}
	return c.conn.WriteMessage(websocket.BinaryMessage, data)
}
