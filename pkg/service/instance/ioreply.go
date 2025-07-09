package instance

import (
	"amoncusir/example/pkg/tools/unres"
	"amoncusir/example/pkg/types"
	"context"
	"errors"
	"io"
	"sync"
)

func newIoReplier(uri *unres.Uri, cnFn ConnFn) *ioReplier {
	return &ioReplier{
		uri:         uri,
		connFn:      cnFn,
		replyBuffer: 32 * 1024, // half of the theorical max for a TCP package size
		replyGroup:  &sync.WaitGroup{},
		cancelLck:   &sync.RWMutex{},
	}
}

type ioReplier struct {
	connFn ConnFn

	uri         *unres.Uri
	replyBuffer int
	replyGroup  *sync.WaitGroup

	cancelLck *sync.RWMutex
}

// Reply implements Instance.
func (r *ioReplier) Reply(ctx context.Context, req types.RequestConn) error {

	r.cancelLck.RLock()

	r.replyGroup.Add(1)
	defer r.replyGroup.Done()

	r.cancelLck.RUnlock()

	return r.reply(ctx, req)
}

func (r *ioReplier) reply(ctx context.Context, req types.RequestConn) error {
	logger.Printf("Reply from %s", req.RemoteAddr())

	dst, err := r.connFn(ctx, r.uri)
	if err != nil {
		return err
	}
	defer dst.Close()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	buf := make([]byte, r.replyBuffer)
	err = ioReply(req, dst, buf)

	return err
}

// Close a current Reply connection
// Do not mislead with Unwatch()
// Cancel wait until all replies are done.
func (r *ioReplier) Close() error {
	r.cancelLck.Lock()
	defer r.cancelLck.Unlock()

	r.replyGroup.Wait()
	return nil
}

// Enable bidirectional communication
func ioReply(source io.ReadWriter, dest io.ReadWriter, buf []byte) error {

	result := make(chan error, 2)
	defer close(result)

	go func() {
		_, err := io.CopyBuffer(dest, source, buf)
		result <- err
	}()

	go func() {
		_, err := io.CopyBuffer(source, dest, buf)
		result <- err
	}()

	r1 := <-result
	r2 := <-result

	return errors.Join(r1, r2)
}
