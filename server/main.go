package main

import (
	"encoding/json"
	"log"
	"net"
	"sync"
)

var (
	controlMu  sync.Mutex
	controlEnc *json.Encoder
)

func main() {
	go listen()
}

// Accepts an incoming request from ngrok lite
func listen() {
	ln, err := net.Listen("tcp", ":4443")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[control] Listening on :4443!")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("[control] Error: ", err)
			continue
		}

		log.Println("[control] Client connected from ", conn.RemoteAddr())

		controlMu.Lock()
		controlEnc = json.NewEncoder(conn)
		controlMu.Unlock()
	}
}
