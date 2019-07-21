package rooms_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/gospeak/chat-engine/rooms"
	"github.com/stretchr/testify/assert"
)

var upgrader = websocket.Upgrader{}

func servWs(w http.ResponseWriter, r *http.Request) {
	_, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func createMockWs() *websocket.Conn {
	s := httptest.NewServer(http.HandlerFunc(servWs))
	defer s.Close()
	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	return ws
}

func TestInterface(t *testing.T) {
	testFun := func(cr rooms.ChatRoom) {}
	cr := rooms.NewGeneralChatRoom()
	testFun(cr)
}

func TestAddClient(t *testing.T) {
	ws := createMockWs()
	defer ws.Close()
	cr := rooms.NewGeneralChatRoom()
	assert.True(t, cr.AddClient(ws, "ID"))
	assert.False(t, cr.AddClient(ws, "ID"), "Add client with same id must be failed")
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, cr.CountClients(), uint32(1))
	assert.True(t, cr.RemoveClientBySessionID("ID"))
	assert.Equal(t, cr.CountClients(), uint32(0))
}

func TestRemoveClient(t *testing.T) {
	ws := createMockWs()
	ws2 := createMockWs()
	defer ws.Close()
	defer ws2.Close()

	cr := rooms.NewGeneralChatRoom()
	assert.True(t, cr.AddClient(ws, "ID"))
	assert.True(t, cr.AddClient(ws2, "ID2"))
	assert.False(t, cr.AddClient(ws, "ID"), "Add client with same id must be failed")
	assert.Equal(t, cr.CountClients(), uint32(2))
	assert.True(t, cr.RemoveClientBySessionID("ID"))
	assert.Equal(t, cr.CountClients(), uint32(1))
	assert.False(t, cr.RemoveClientBySessionID("ID"))
	assert.Equal(t, cr.CountClients(), uint32(1))
}
