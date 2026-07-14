package client

import (
	"encoding/json"
	"log"
	"net"
	"nsky/internal/protocol"
)

func main() {
	conn, err := net.Dial("tcp", "3000:4443")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[client] Connecting to control channel in 3000:4443")
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

}
