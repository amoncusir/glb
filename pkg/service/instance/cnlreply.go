package instance

import (
	"amoncusir/example/pkg/tools/ioutils"
	"amoncusir/example/pkg/tools/unres"
	"amoncusir/example/pkg/types"
	"context"
	"sync"
)

func newCancelableReplier(uri *unres.Uri, cnFn ConnFn) *cancelableReplier {
	return &cancelableReplier{
		uri:         uri,
		connFn:      cnFn,
		replyBuffer: 32 * 1024, // half of the theorical max for a TCP package size
		replyGroup:  &sync.WaitGroup{},
		cancelLck:   &sync.RWMutex{},
	}
}

type cancelableReplier struct {
	connFn ConnFn

	uri         *unres.Uri
	replyBuffer int
	replyGroup  *sync.WaitGroup

	cancelLck *sync.RWMutex
}

// Reply implements Instance.
func (r *cancelableReplier) Reply(ctx context.Context, req types.RequestConn) error {

	r.cancelLck.RLock()

	r.replyGroup.Add(1)
	defer r.replyGroup.Done()

	r.cancelLck.RUnlock()

	return r.reply(ctx, req)
}

func (r *cancelableReplier) reply(ctx context.Context, req types.RequestConn) error {
	logger.Printf("Reply from %s", req.RemoteAddr())

	dst, err := r.connFn(ctx, r.uri)
	if err != nil {
		return err
	}

	defer dst.Close()

	buf := make([]byte, r.replyBuffer)
	err = ioutils.Reply(ctx, req, dst, buf)

	return err
}

// Close a current Reply connection
// Do not mislead with Unwatch()
// Cancel wait until all replies are done.
func (r *cancelableReplier) Close() error {
	r.cancelLck.Lock()
	defer r.cancelLck.Unlock()

	r.replyGroup.Wait()
	return nil
}
