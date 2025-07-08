package ioutils

import (
	"context"
	"errors"
	"io"
)

func CtxCopy(ctx context.Context, source io.Reader, dest io.Writer, buf []byte) (written int64, err error) {
	for {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}

		nr, er := source.Read(buf)

		if nr > 0 {
			nw, ew := dest.Write(buf[0:nr])

			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					err = errors.New("invalid write result")
					break
				}
			}

			written += int64(nw)

			if ew != nil {
				err = ew
				break
			}

			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}

		if er != nil {
			if er != io.EOF {
				err = er
				break
			}
			break
		}
	}

	return written, err
}
