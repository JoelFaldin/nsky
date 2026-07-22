package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"nsky/internal/protocol"
	"sync"
)

var (
	controlMu  sync.Mutex
	controlEnc *json.Encoder

	pendingMu sync.Mutex
	pending   = map[string]net.Conn{}
)

func main() {
	go listenControl()
	go listenJoin()
	listenPublic()
}

// Accepts an incoming request from ngrok lite
func listenControl() {
	ln, err := net.Listen("tcp", ":4443")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[control] Listening on :4443!")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("[control] Acept error: ", err)
			continue
		}

		log.Println("[control] Client connected from ", conn.RemoteAddr())

		controlMu.Lock()
		controlEnc = json.NewEncoder(conn)
		controlMu.Unlock()
	}
}

// Accepts an internet visitor. For each one, ask the client (control)
// to open a join connection
func listenPublic() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[public] Listening on :8080!")

	var counter int
	for {
		visitor, err := ln.Accept()
		if err != nil {
			log.Println("[public] Accept error: ", err)
			continue
		}

		controlMu.Lock()
		enc := controlEnc
		controlMu.Unlock()

		if enc == nil {
			log.Println("[public] No client connected, rejecing visitor...")
			visitor.Close()
			continue
		}

		counter++
		id := fmt.Sprintf("%d", counter)

		pendingMu.Lock()
		pending[id] = visitor
		pendingMu.Unlock()

		log.Println("[public] Visitor", visitor.RemoteAddr(), "-> id", id)

		controlMu.Lock()
		err = enc.Encode(protocol.Message{Type: "new_conn", ID: id})
		controlMu.Unlock()
		if err != nil {
			log.Println("[control] Couldnt report to client:", err)
		}
	}
}

func listenJoin() {
	ln, err := net.Listen("tcp", ":4444")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[join] Listening on :4444!")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("[join] Accept error:", err)
			continue
		}

		go handleJoin(conn)
	}
}

func handleJoin(conn net.Conn) {
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Println("[handle] Error when reading:", err)
	}

	fmt.Printf("[handle] Reading %q\n", line)
}
