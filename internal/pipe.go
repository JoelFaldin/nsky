package internal

import (
	"io"
	"net"
)

func Pipe(a, b net.Conn) {
	channel := make(chan struct{}, 2)

	go func() { io.Copy(a, b); channel <- struct{}{} }()
	go func() { io.Copy(b, a); channel <- struct{}{} }()
	<-channel

	a.Close()
	b.Close()
}
