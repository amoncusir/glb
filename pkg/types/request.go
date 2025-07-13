package types

import (
	"net"
)

type RequestConn interface {
	net.Conn

	RemoteIp() string
}
