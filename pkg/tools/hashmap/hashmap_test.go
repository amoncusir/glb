package hashmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasmapSetAndGet(t *testing.T) {
	assert := assert.New(t)

	m := New(16, 5./16., 1.5)

	m.Set([]byte{0x00, 0xFF, 0x0F}, 1)
	m.Set([]byte{0x0A, 0xFA, 0x3F}, 2)
	m.Set([]byte{0x5D, 0x9E, 0x1F}, 3)

	assert.Equal(1, m.Get([]byte{0x00, 0xFF, 0x0F}))
	assert.Equal(2, m.Get([]byte{0x0A, 0xFA, 0x3F}))
	assert.Equal(3, m.Get([]byte{0x5D, 0x9E, 0x1F}))
}
