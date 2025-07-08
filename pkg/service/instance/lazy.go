package instance

import (
	"amoncusir/example/pkg/tools/unres"
	"context"
	"errors"
	"sync/atomic"
	"time"
)

func newLazyHealthcheck(uri *unres.Uri, cnFn ConnFn) *lazyHealthcheck {
	return &lazyHealthcheck{
		uri:    uri,
		connFn: cnFn,

		Interval: 5 * time.Second,

		CallbackHealthy:   func() {},
		CallbackUnhealthy: func() {},

		watching:      &atomic.Bool{},
		statusChannel: make(chan Status),
		lastStatus:    &atomic.Int64{},
	}
}

type lazyHealthcheck struct {
	_ noCopy

	uri    *unres.Uri
	connFn ConnFn

	Interval time.Duration
	Timeout  time.Duration

	CallbackHealthy   func()
	CallbackUnhealthy func()

	watching *atomic.Bool

	statusChannel chan Status
	lastStatus    *atomic.Int64
}

// Status implements Instance.
func (lazy *lazyHealthcheck) Status() Status {
	return Status(lazy.lastStatus.Load())
}

func (lazy *lazyHealthcheck) setStatus(status Status) {
	lazy.lastStatus.Store(int64(status))

	if status == STATUS_HEALTHY {
		if lazy.CallbackHealthy != nil {
			go lazy.CallbackHealthy()
		}
	} else if lazy.CallbackUnhealthy != nil {
		go lazy.CallbackUnhealthy()
	}
}

func (lazy *lazyHealthcheck) SetHealthy() error {
	return lazy.setHealth(STATUS_HEALTHY)
}

func (lazy *lazyHealthcheck) SetUnhealthy() error {
	return lazy.setHealth(STATUS_UNHEALTHY)
}

func (lazy *lazyHealthcheck) setHealth(status Status) error {
	if !lazy.watching.Load() {
		return errors.New("must be watching to set health")
	}

	lazy.statusChannel <- status

	return nil
}

// Watch implements Instance.
func (lazy *lazyHealthcheck) Watch() error {

	if !lazy.watching.CompareAndSwap(false, true) {
		return errors.New("already watching")
	}

	defer func() {
		lazy.watching.Store(false)
	}()

	lazy.watch()

	return nil
}

// Only one goroutine must run this code
func (lazy *lazyHealthcheck) watch() {

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	for status := range lazy.statusChannel {
		logger.Printf("New instance status update: %d", status)

		switch status {
		case STATUS_HEALTHY:
			if lazy.lastStatus.Swap(int64(status)) != int64(STATUS_HEALTHY) {
				lazy.setStatus(status)
			}
		case STATUS_UNHEALTHY:
			if lazy.lastStatus.Swap(int64(status)) != int64(STATUS_UNHEALTHY) {
				lazy.setStatus(status)
				go lazy.watchUntilHealthy(ctx)
			}
		case STATUS_UNKNOWN:
			// Unwatch called.
			return
		}
	}
}

func (lazy *lazyHealthcheck) watchUntilHealthy(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if lazy.Timeout > 0 {
			c, ctxClose := context.WithTimeout(ctx, lazy.Timeout)
			defer ctxClose()
			ctx = c
		}

		conn, err := lazy.connFn(ctx, lazy.uri)

		if err == nil {
			conn.Close()
			break
		}

		time.Sleep(lazy.Interval)
	}
}

// Watch implements Instance.
func (lazy *lazyHealthcheck) Unwatch() error {

	if !lazy.watching.Load() {
		return errors.New("not watching")
	}

	lazy.statusChannel <- STATUS_UNKNOWN

	return nil
}
