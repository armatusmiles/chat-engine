package main

import (
	"flag"
	"net/http"

	"github.com/gospeak/chat-engine/api"
	log "github.com/sirupsen/logrus"
)

var addr = flag.String("addr", ":8888", "http service address")
var debugLevel = flag.Int("debugLevel", 5, "")

func main() {
	log.SetLevel(log.Level(*debugLevel))
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	flag.Parse()

	router := api.Init()

	http.Handle("/", router)
	log.Info("Server listening on ", *addr)

	log.Fatal("ListenAndServe:", http.ListenAndServe(*addr, nil))
}
