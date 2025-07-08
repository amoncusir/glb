package ioutils

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func newRW(name string, chunks ...string) *RWChunked {
	inst := &RWChunked{
		name: name,
	}

	for i, c := range chunks {
		if i == 0 {
			inst.rbuffer = []byte(c)
		} else {
			inst.rchunks = append(inst.rchunks, []byte(c))
		}
	}

	return inst
}

type RWChunked struct {
	name    string
	rbuffer []byte
	rchunks [][]byte

	wchunks [][]byte
}

func (rw *RWChunked) Read(p []byte) (n int, err error) {
	bl := len(rw.rbuffer)

	if bl <= 0 {
		return 0, io.EOF
	}

	pl := len(p)

	if pl <= 0 {
		return 0, io.ErrShortBuffer
	}

	n = min(pl, bl)

	copy(p, rw.rbuffer[:n])

	rw.rbuffer = rw.rbuffer[n:]

	if len(rw.rbuffer) <= 0 {
		err = io.EOF

		if len(rw.rchunks) > 0 {
			rw.rbuffer = rw.rchunks[0]

			if len(rw.rchunks) > 1 {
				rw.rchunks = rw.rchunks[1:]
			} else {
				rw.rchunks = make([][]byte, 0)
			}
		}
	}

	return n, err
}

func (rw *RWChunked) Write(p []byte) (n int, err error) {
	dst := make([]byte, len(p))
	copy(dst, p)

	rw.wchunks = append(rw.wchunks, dst)

	return len(p), nil
}

func (rw *RWChunked) WChunks() []string {
	s := make([]string, len(rw.wchunks))

	for i, v := range rw.wchunks {
		s[i] = string(v[:])
	}

	return s
}

func TestReplyFnWithLargeBuffer(t *testing.T) {
	assert := assert.New(t)

	source := newRW("source", "first", "third", "fifth")
	dest := newRW("dest", "second", "fourth")
	buff := make([]byte, 64)

	// Act
	Reply(context.Background(), (io.ReadWriter)(source), (io.ReadWriter)(dest), buff)

	assert.Equal([]string{"second", "fourth"}, source.WChunks())
	assert.Equal([]string{"first", "third", "fifth"}, dest.WChunks())
}

func TestReplyFnWithSmallBuffer(t *testing.T) {
	assert := assert.New(t)

	source := newRW("source", "first", "third", "fifth")
	dest := newRW("dest", "second", "fourth")
	buff := make([]byte, 1)

	// Act
	Reply(context.Background(), (io.ReadWriter)(source), (io.ReadWriter)(dest), buff)

	assert.Equal([]string{"s", "e", "c", "o", "n", "d", "f", "o", "u", "r", "t", "h"}, source.WChunks())
	assert.Equal([]string{"f", "i", "r", "s", "t", "t", "h", "i", "r", "d", "f", "i", "f", "t", "h"}, dest.WChunks())
}
