package filter

import (
	"amoncusir/example/pkg/tools/bitmapcheck"
	"amoncusir/example/pkg/types"
	"context"
	"errors"
	"net"
)

func NewIpFilter(allow bool, ips ...net.IP) PreFilter {

	var ips4, ips6 []net.IP

	for _, ip := range ips {
		if len(ip) == net.IPv4len {
			ips4 = append(ips4, ip)
		} else if len(ip) == net.IPv6len {
			ips6 = append(ips6, ip)
		}
	}

	return &ipAllower{
		allow: allow,
		ip4:   generateIp4(ips4),
		ip6:   generateIp6(ips6),
	}
}

func generateIp4(ips []net.IP) [net.IPv4len]bitmapcheck.BinBitmap {
	panic("TODO")
}

func generateIp6(ips []net.IP) [net.IPv6len]bitmapcheck.BinBitmap {
	panic("TODO")
}

type ipAllower struct {
	allow bool
	ip4   [net.IPv4len]bitmapcheck.BinBitmap
	ip6   [net.IPv6len]bitmapcheck.BinBitmap
}

func (r *ipAllower) Accept(ctx context.Context, req types.RequestConn) error {
	ip := req.RemoteIp()

	if len(ip) == net.IPv4len {
		if r.acceptIp4(ip) {
			return nil
		}
	} else if len(ip) == net.IPv6len {
		if r.acceptIp6(ip) {
			return nil
		}
	}

	return errors.New("filtered IP")
}

func (r *ipAllower) acceptIp4(ip net.IP) bool {
	for i, c := range r.ip4 {
		if !c.Get(int(ip[i])) {
			return !r.allow
		}
	}

	return r.allow
}

func (r *ipAllower) acceptIp6(ip net.IP) bool {
	for i, c := range r.ip6 {
		if !c.Get(int(ip[i])) {
			return !r.allow
		}
	}

	return r.allow
}
