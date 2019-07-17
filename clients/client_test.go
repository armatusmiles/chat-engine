package clients_test

import (
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/gospeak/chat-engine/clients"
	"github.com/gospeak/protorepo/dbmanage"
	"github.com/stretchr/testify/assert"
)

var mockMsg = dbmanage.ChatMessage{UserId: "111", Message: "Hello client"}

func servWs(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	mockChan := make(chan dbmanage.ChatMessage)
	chatClient := clients.NewChatClient(c, "mock-session-id", mockChan)

	chatClient.SendMessage(&mockMsg)
	chatClient.CloseConnection()
}

func TestCloseConnection(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(servWs))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	for i := 0; i < 50; i++ {
		_, _, err := websocket.DefaultDialer.Dial(u, nil)
		assert.Nil(t, err)
	}
	assert.True(t, runtime.NumGoroutine() < 10, "Memory leak gorutine detected")
}

func TestSendMessage(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(servWs))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	assert.Nil(t, err)
	defer ws.Close()

	// wait from mock message from server
	mt, p, err := ws.ReadMessage()
	assert.Nil(t, err)
	assert.Equal(t, mt, websocket.BinaryMessage, "Message type must be binary")

	msg := dbmanage.ChatMessage{}
	err = proto.Unmarshal(p, &msg)
	assert.Nil(t, err)

	assert.Equal(t, msg.UserId, mockMsg.UserId, "Wrong userID response")
	assert.Equal(t, msg.Message, mockMsg.Message, "Wrong Message response")
}

func TestSendNotProtoMsg(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(servWs))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	assert.Nil(t, err)

	chatClient := clients.NewChatClient(ws, "mock-session-id", make(chan dbmanage.ChatMessage))
	err = chatClient.SendMessage(nil)
	assert.NotNil(t, err)
}
