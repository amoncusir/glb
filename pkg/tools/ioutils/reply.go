package ioutils

import (
	"context"
	"errors"
	"io"
)

type replyResult struct {
	w   int64
	err error
}

// Enable bidirectional communication
func Reply(ctx context.Context, source io.ReadWriter, dest io.ReadWriter, buf []byte) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	result := make(chan replyResult, 2)
	defer close(result)

	go func() {
		w, err := CtxCopy(ctx, dest, source, buf)
		result <- replyResult{
			w:   w,
			err: err,
		}
	}()

	go func() {
		w, err := CtxCopy(ctx, source, dest, buf)
		result <- replyResult{
			w:   w,
			err: err,
		}
	}()

	r1 := <-result
	r2 := <-result

	return errors.Join(r1.err, r2.err)
}
