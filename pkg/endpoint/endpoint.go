package endpoint

import (
	"io"

	"amoncusir/example/pkg/router"
)

const (
	PROTOCOL_TCP string = "tcp"
	PROTOCOL_UDP string = "udp"
)

const (
	STATUS_LISTEN int = 1
	STATUS_CLOSED int = 0
)

type Endpoint interface {
	Listen(rt router.Router) error
	Status() <-chan int

	Protocol() string
	Address() string

	io.Closer
}
