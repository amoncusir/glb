package endpoint

import (
	"amoncusir/example/pkg/router"
	"amoncusir/example/pkg/tools/unres"
	"amoncusir/example/pkg/types"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var (
	logger = log.New(os.Stdout, "endpoint/net_socket:", log.LstdFlags)
)

func NewSocket(uri *unres.Uri) (Endpoint, error) {
	if uri.Scheme == "" || uri.Autority.String() == "" {
		return nil, errors.New("invalid uri must follow the format <protocol>://<host>:<port>")
	}

	return &socketEndpoint{
		protocol: uri.Scheme,
		address:  uri.Autority.String(),

		statusListeners: []chan int{},
		listener:        nil,
	}, nil
}

type socketEndpoint struct {
	address         string
	protocol        string
	timeout         time.Duration
	statusListeners []chan int
	listener        net.Listener
}

func (e *socketEndpoint) String() string {
	return fmt.Sprintf("SocketEndpoint{%s://%s, timeout=%d}", e.protocol, e.address, e.timeout)
}

func (e *socketEndpoint) Protocol() string {
	return e.protocol
}

func (e *socketEndpoint) Address() string {
	return e.address
}

func (e *socketEndpoint) Close() error {
	ln := e.listener

	if ln == nil {
		return errors.New("no socket to close")
	}

	defer func() {
		e.listener = nil
		e.setStatus(STATUS_CLOSED)
	}()

	return ln.Close()
}

func (e *socketEndpoint) Status() <-chan int {
	nw := make(chan int)

	e.statusListeners = append(e.statusListeners, nw)

	return nw
}
func (e *socketEndpoint) setStatus(status int) {
	for _, c := range e.statusListeners {
		go func() {
			c <- status
		}()
	}
}

func (e *socketEndpoint) Listen(rt router.Router) error {
	if rt == nil {
		return errors.New("router can be nil")
	}

	if e.listener != nil {
		return errors.New("already opened connection")
	}

	if ln, err := net.Listen(string(e.protocol), e.address); err != nil {
		return err
	} else {
		e.listener = ln
	}

	go e.listenLoop(rt)
	return nil
}

func (e *socketEndpoint) listenLoop(rt router.Router) {
	defer e.Close()

	e.setStatus(STATUS_LISTEN)

	for {
		conn, err := e.listener.Accept()

		if err != nil {
			logger.Printf("Error when retrieve data from %s. %s", e, err)
			break
		}

		go func() {
			ctx := context.Background()
			ctx, c := context.WithTimeout(ctx, e.timeout)
			err := rt.RouteRequest(ctx, types.RequestConn(conn))

			defer conn.Close()
			defer c()

			if err != nil {
				logger.Print("Error on route: ", err)
			}
		}()
	}
}
