package endpoint

import (
	"amoncusir/example/pkg/types"
	"net"
	"strings"
)

func reqFromConn(conn net.Conn) types.RequestConn {
	return &req{
		Conn: conn,
	}
}

type req struct {
	net.Conn
}

// RemoteIp implements types.RequestConn.
func (r *req) RemoteIp() string {
	rAdd := r.RemoteAddr().String()

	is6 := strings.HasPrefix(rAdd, "[")

	if is6 {
		if strings.HasSuffix(rAdd, "]") { // No port specified
			return rAdd
		}

		li := strings.LastIndex(rAdd, ":")
		return rAdd[:li]
	} else {
		return strings.Split(rAdd, ":")[0]
	}
}
