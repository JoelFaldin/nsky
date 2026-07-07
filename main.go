package main

import (
	"log"
	"time"
)

func main() {
	for {
		if err := connectAndServeOnce(); err != nil {
			log.Println("session ended:", err)
		}

		log.Println("Reconnecting...")
		time.Sleep(1 * time.Second)
	}
}
