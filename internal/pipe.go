package internal

import (
	"io"
	"net"
)

func Pipe(join, local net.Conn) {
	go func() { io.Copy(join, local) }
	go func() { io.Copy(local, join) }
}
