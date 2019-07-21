package clients_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/gospeak/chat-engine/clients"
	dbm "github.com/gospeak/protorepo/dbmanage"
	"github.com/stretchr/testify/assert"
)

// expected count after call initMockClientsAndServ
const mockClientsExpectedCount = 4

// Don't forget close a server!
func initMockClientsAndServ() (*httptest.Server, *clients.ClientList) {
	dummyServeWs := func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{}
		upgrader.Upgrade(w, r, nil)
	}
	s := httptest.NewServer(http.HandlerFunc(dummyServeWs))

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Create websocket.Conn for mockClients
	var wsCLients [mockClientsExpectedCount]*websocket.Conn
	for i := 0; i < mockClientsExpectedCount; i++ {
		ws, _, err := websocket.DefaultDialer.Dial(u, nil)
		if err != nil {
			panic(err)
		}
		wsCLients[i] = ws
	}

	// Don't change sessionID's and count of clients They are use in another tests
	chatClients := []clients.ChatClient{
		*clients.NewChatClient(wsCLients[0], "1", make(chan dbm.ChatMessage)),
		*clients.NewChatClient(wsCLients[1], "2", make(chan dbm.ChatMessage)),
		*clients.NewChatClient(wsCLients[2], "3", make(chan dbm.ChatMessage)),
		*clients.NewChatClient(wsCLients[3], "4", make(chan dbm.ChatMessage)),
	}
	cl := clients.NewChatClientList()

	for _, v := range chatClients {
		if cl.Add(v) != true {
			panic("Error add Chat Client")
		}
	}
	if cl.Count() != uint32(len(chatClients)) {
		panic("Count clients is wrong")
	}
	return s, cl
}

func TestNewClientList(t *testing.T) {
	cl := clients.NewChatClientList()
	assert.NotEqual(t, cl, nil, "List should't be nil")
	assert.Zero(t, cl.Count())
}

func TestAddClient(t *testing.T) {

	s, cl := initMockClientsAndServ()
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	assert.Nil(t, err)
	client := clients.NewChatClient(ws, "999", make(chan dbm.ChatMessage))

	clientsCountBeforeAdd := cl.Count()

	assert.True(t, cl.Add(*client))
	assert.Equal(t, cl.Count(), clientsCountBeforeAdd+1)

	// Try add duplicate
	assert.False(t, cl.Add(*client))
	assert.Equal(t, cl.Count(), clientsCountBeforeAdd+1)
}

func TestRemoveClient(t *testing.T) {
	s, cl := initMockClientsAndServ()
	defer s.Close()
	countBefore := cl.Count()
	assert.True(t, cl.RemoveBySessionID("1"))
	assert.Equal(t, countBefore-1, cl.Count())

	// try remove not exists client
	assert.False(t, cl.RemoveBySessionID("1"))
	assert.Equal(t, countBefore-1, cl.Count())
}

func TestIsExists(t *testing.T) {
	s, cl := initMockClientsAndServ()
	defer s.Close()
	assert.True(t, cl.IsExists("1"))
	assert.False(t, cl.IsExists("notExistsKey"))
}
