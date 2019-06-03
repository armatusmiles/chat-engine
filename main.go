package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()
var addr = flag.String("addr", ":8888", "http service address")
var debugLevel = flag.Int("debugLevel", 5, "")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Info(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	// TODO
	// Add check session-id. If id is invalid shutdown connection

	defer conn.Close()
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Debugln(err)
			break
		}
		messageHandler(conn, mt, message)
		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Debugln(err)
			break
		}
	}
}

func messageHandler(con *websocket.Conn, messageType int, message []byte) {
	switch messageType {
	case websocket.TextMessage:
	case websocket.BinaryMessage:
		log.Warning("Binary message not supported for now. Connection is closed")
		con.Close()
	case websocket.CloseMessage:
		con.Close()
	case websocket.PingMessage:
		// TODO send pong?
	case websocket.PongMessage:
		// Do nothing
	default:
		log.Warning("Unknown websocket message code. Connection is closed")
	}
}

func main() {
	log.Out = os.Stdout
	log.SetLevel(logrus.Level(*debugLevel))
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	flag.Parse()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})

	log.Info("Server listening on ", *addr)

	err := http.ListenAndServe(*addr, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
