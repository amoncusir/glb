package instance

import (
	"amoncusir/example/pkg/tools/unres"
	"amoncusir/example/pkg/types"
	"context"
	"errors"
	"net"
)

const (
	REPLY_BUFFER int = 32 * 1024
)

type ConnFn func(ctx context.Context, uri *unres.Uri) (net.Conn, error)

func New(uri *unres.Uri) Instance {
	connFn := func(ctx context.Context, _ *unres.Uri) (net.Conn, error) {
		d := &net.Dialer{}
		return d.DialContext(ctx, uri.Scheme, uri.Autority.String())
	}

	replier := newCancelableReplier(uri, connFn)
	healthcheck := newLazyHealthcheck(uri, connFn)

	return &lazyInstance{
		cancelableReplier: replier,
		lazyHealthcheck:   healthcheck,
	}
}

type lazyInstance struct {
	_ noCopy

	*lazyHealthcheck
	*cancelableReplier
}

// AddHealthyCallback implements Instance.
func (i *lazyInstance) AddHealthyCallback(fn func(self Instance)) error {
	if i.lazyHealthcheck.CallbackHealthy != nil {
		return errors.New("exist callback")
	}

	i.lazyHealthcheck.CallbackHealthy = func() { fn(i) }
	return nil
}

// AddUnhealthyCallback implements Instance.
func (i *lazyInstance) AddUnhealthyCallback(fn func(self Instance)) error {
	if i.lazyHealthcheck.CallbackUnhealthy != nil {
		return errors.New("exist callback")
	}

	i.lazyHealthcheck.CallbackUnhealthy = func() { fn(i) }
	return nil
}

// Address implements Instance.
func (i *lazyInstance) Protocol() string {
	return i.cancelableReplier.uri.Scheme
}

// Address implements Instance.
func (i *lazyInstance) Address() string {
	return i.cancelableReplier.uri.Autority.String()
}

// Reply implements Instance.
func (i *lazyInstance) Reply(ctx context.Context, req types.RequestConn) error {
	err := i.cancelableReplier.Reply(ctx, req)

	if err != nil {
		i.lazyHealthcheck.SetUnhealthy()
	}

	return err
}
