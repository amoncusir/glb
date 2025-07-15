package bitmapcheck

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinBitmap(t *testing.T) {
	assert := assert.New(t)
	mapCheck := []bool{
		true, false, false, true, false, false, false, false, true,
	}

	bm := NewBoolean(len(mapCheck))

	for i, v := range mapCheck {
		bm.Set(i, v)
	}

	t.Log(bm.Size())
	t.Log(len(mapCheck))

	for i, v := range mapCheck {
		assert.Equal(v, bm.Get(i))
	}
}
