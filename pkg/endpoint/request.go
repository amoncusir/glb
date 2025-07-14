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
func (r *req) RemoteIp() net.IP {
	var ip string

	rAdd := r.RemoteAddr().String()
	is6 := strings.HasPrefix(rAdd, "[")

	if is6 {
		if !strings.HasSuffix(rAdd, "]") { // Port specified
			li := strings.LastIndex(rAdd, ":")
			rAdd = rAdd[:li]
		}

		ip = rAdd[1 : len(rAdd)-1]
	} else {
		ip = strings.Split(rAdd, ":")[0]
	}

	return net.ParseIP(ip)
}

func (r *req) RemoteIpString() string {
	return r.RemoteIp().String()
}
