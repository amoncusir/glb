package ioutils

import (
	"context"
	"io"
)

// Reply bytes in both ways until any RW has more content
func Reply(ctx context.Context, source io.ReadWriter, dest io.ReadWriter, buf []byte) error {

	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		s, d := source, dest
		if i%2 != 0 {
			// Dest -> Source
			s = dest
			d = source
		}

		if w, err := CtxCopy(ctx, s, d, buf); err != nil {
			return err
		} else if w <= 0 {
			// No sended bytes
			break
		}
	}

	return nil
}
