package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()
var addr = flag.String("addr", ":8080", "http service address")
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
	defer conn.Close()
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Debug("Read:", err)
			break
		}
		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Debug("Write:", err)
			break
		}
	}
}

func main() {
	log.Out = os.Stdout
	log.SetLevel(logrus.Level(*debugLevel))

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
