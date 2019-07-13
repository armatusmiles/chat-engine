package broadcast_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/gospeak/chat-engine/clients"
	"github.com/gospeak/chat-engine/rooms/broadcast"
	"github.com/gospeak/protorepo/dbmanage"
)

var upgrader = websocket.Upgrader{}
var mockClients []clients.ChatClient
var broadcaster *broadcast.Broadcaster

const mockClientsExpectedCount = 3

var mockMsg = dbmanage.ChatMessage{UserId: "111", Message: "Hello client"}
var serveMutex = sync.Mutex{}

func servWs(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	serveMutex.Lock()
	chatClient := clients.NewChatClient(c, "mock-session-id", broadcaster.ReadCh)
	mockClients = append(mockClients, *chatClient)
	broadcaster.UpdateListClients(&mockClients)
	if len(mockClients) == mockClientsExpectedCount {
		serveMutex.Unlock()
		broadcaster.BroadcastMessage(&mockMsg)
		beforeClose := runtime.NumGoroutine()
		broadcaster.Close()
		afterClose := runtime.NumGoroutine()
		if beforeClose-1 != afterClose {
			panic("broadcaster Close is failed. chReader gorutine not been stopped")
		}
		return
	}
	serveMutex.Unlock()
}

func TestBroadcast(t *testing.T) {
	broadcaster = broadcast.NewBroadcaster(&mockClients)
	s := httptest.NewServer(http.HandlerFunc(servWs))
	defer s.Close()
	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	var wsCLients [mockClientsExpectedCount]*websocket.Conn
	for i := 0; i < mockClientsExpectedCount; i++ {
		ws, _, err := websocket.DefaultDialer.Dial(u, nil)
		if err != nil {
			t.Fatalf("%v", err)
		}
		wsCLients[i] = ws
	}
	for i := 0; i < mockClientsExpectedCount; i++ {
		var msg dbmanage.ChatMessage
		mt, bMsg, err := wsCLients[i].ReadMessage()
		if err != nil {
			t.Fatal(err)
		}
		err = proto.Unmarshal(bMsg, &msg)
		if err != nil {
			t.Fatal(err)
		}
		if msg.Message != mockMsg.Message {
			t.Fatal("returned bad text message")
		}

		if msg.UserId != mockMsg.UserId {
			t.Fatal("returned bad ID sender message")
		}

		if mt != websocket.BinaryMessage {
			t.Fatal("Message type must be Binary message")
		}
	}
}
