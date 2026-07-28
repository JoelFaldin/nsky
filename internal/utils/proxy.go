package utils

import (
	"io"
	"net"
)

func Proxy(visitor, join net.Conn) {
	done := make(chan struct{})
	go func() {
		io.Copy(visitor, join)
		close(done)
	}()

	io.Copy(join, visitor)
	<-done
}
