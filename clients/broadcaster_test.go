package clients_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/gospeak/chat-engine/clients"
	"github.com/gospeak/protorepo/dbmanage"
	"github.com/stretchr/testify/assert"
)

func TestBroadcast(t *testing.T) {

	var upgrader = websocket.Upgrader{}
	clientList := clients.NewChatClientList()
	broadcaster := clients.NewBroadcaster(clientList)
	const mockClientsExpectedCount = 3
	var id = 0
	mux := &sync.Mutex{}

	servWs := func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
		}

		mux.Lock() // mux for fix data race
		chatClient := clients.NewChatClient(c, strconv.Itoa(id), broadcaster.ReadCh)
		id++
		mux.Unlock()
		clientList.Add(*chatClient)

		// When count of clients will be mockClientsExpectedCount send them mock message
		if clientList.Count() == mockClientsExpectedCount {
			broadcaster.BroadcastMessage(&mockMsg)
			broadcaster.Close()
		}
	}

	s := httptest.NewServer(http.HandlerFunc(servWs))
	defer s.Close()
	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Create websocket.Conn for mockClients
	var wsCLients [mockClientsExpectedCount]*websocket.Conn
	for i := 0; i < mockClientsExpectedCount; i++ {
		ws, _, err := websocket.DefaultDialer.Dial(u, nil)
		assert.Nil(t, err)
		wsCLients[i] = ws
	}

	for i := 0; i < mockClientsExpectedCount; i++ {
		var msg dbmanage.ChatMessage
		// block here and wait for broadcast message from server
		// all clients must should get a message
		mt, bMsg, err := wsCLients[i].ReadMessage()
		assert.Nil(t, err)
		err = proto.Unmarshal(bMsg, &msg)
		assert.Nil(t, err)
		assert.Equal(t, msg.Message, mockMsg.Message, "response message must be the same of mock")
		assert.Equal(t, msg.UserId, mockMsg.UserId, "bad ID sender message")
		assert.Equal(t, mt, websocket.BinaryMessage, "Message type must be binary")
	}
}
