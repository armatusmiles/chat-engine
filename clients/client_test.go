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
)

var upgrader = websocket.Upgrader{}

var mockMsg = dbmanage.ChatMessage{UserId: "111", Message: "Hello client"}

func servWs(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	mockChan := make(chan dbmanage.ChatMessage)
	chatClient := clients.NewChatClient(c, "mock-session-id", mockChan)

	chatClient.SendMessage(&mockMsg)
	c.Close()
}

func servCloseConnTest(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	mockChan := make(chan dbmanage.ChatMessage)
	chatClient := clients.NewChatClient(c, "mock-session-id", mockChan)
	chatClient.CloseConnection()
}

func TestCloseConnection(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(servWs))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	for i := 0; i < 50; i++ {
		_, _, err := websocket.DefaultDialer.Dial(u, nil)
		if err != nil {
			t.Fatalf("%v", err)
		}
	}
	if runtime.NumGoroutine() > 10 {
		t.Fatal("Memory leak gorutine detected")
	}
}

func TestSendMessage(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(servWs))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	mt, p, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	if mt != websocket.BinaryMessage {
		t.Fatal("Type of message must be binary (proto)")
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
