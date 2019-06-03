package clients_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/armatusmiles/gospeak-chat-engine/clients"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/gospeak/protorepo/dbmanage"
)

var upgrader = websocket.Upgrader{}

var mockMsg = dbmanage.ChatMessage{UserId: 111, Message: "Hello client"}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	chatClient := clients.NewChatClient(c, "mock-session-id")

	_, _, err = c.ReadMessage()
	if err != nil {
		return
	}

	chatClient.SendMessage(&mockMsg)
	c.Close()
}

func TestSendMessage(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	if err := ws.WriteMessage(websocket.TextMessage, []byte("someProtoData")); err != nil {
		t.Fatalf("%v", err)
	}

	_, p, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	msg := dbmanage.ChatMessage{}
	proto.Unmarshal(p, &msg)

	if msg.UserId != mockMsg.UserId {
		t.Fatal("Wrong userID response")
	}

	if msg.Message != mockMsg.Message {
		t.Fatal("Wrong Message response")
	}
}
