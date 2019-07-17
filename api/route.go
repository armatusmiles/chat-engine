package api

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/gospeak/chat-engine/rooms"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// temporary for testing!
var mockRoom = rooms.NewGeneralChatRoom()

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveHome(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "home.html")
}

func serveWs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	sessionID, err := r.Cookie("session_id")
	if err != nil {
		log.Warn("Request without session_id in cookie")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Debugf("New client connection from: %s", r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}

	mockRoom.AddClient(conn, sessionID.Value)
}

func Init() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", serveHome)
	router.GET("/rooms/:id/enter", serveWs)
	return router
}
