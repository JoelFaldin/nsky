package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"nsky/internal/protocol"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:4443")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[client] Connecting to control channel in localhost:4443")
	log.Println("[client] Forwarding traffic to :3000")

	dec := json.NewDecoder(conn)
	for {
		var msg protocol.Message
		if err := dec.Decode(&msg); err != nil {
			log.Fatal("[client] Control conection closed:", err)
		}

		switch msg.Type {
		case "new_conn":
			go handleNewConn(msg.ID, "3000:4444", "3000")
		default:
			log.Println("[client] Unknown message:", msg.Type)
		}
	}
}

func handleNewConn(id, join, local string) {
	joinConn, err := net.Dial("tcp", join)
	if err != nil {
		log.Fatalf("Couldnt open join conn: %v", err)
	}

	fmt.Fprintf(joinConn, "%s\n", id)

	localConn, err := net.Dial("tcp", local)
	if err != nil {
		log.Println("Couldnt connect to local app: ", err)
		joinConn.Close()
		return
	}

	log.Println("Stream ", id, " -> proxing towards ", local)
	pipe(joinConn, localConn)

}
