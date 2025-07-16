package hashmap

import (
	v2 "math/rand/v2"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[v2.IntN(len(letterRunes))]
	}
	return string(b)
}

type P struct {
	k string
	v int
}

func TestHasmapMetadataSetAndGet(t *testing.T) {
	assert := assert.New(t)

	m := NewM8(200, .8/10., 2.)

	pairs := make([]*P, 1_000_000)

	for i := range pairs {
		pairs[i] = &P{
			k: strconv.Itoa(i),
			v: v2.Int(),
		}
	}

	for _, p := range pairs {
		m.Set(p.k, p.v)
	}

	for _, p := range pairs {
		assert.Equal(p.v, m.Get(p.k))
	}
}

func TestDefaultSetAndGet(t *testing.T) {
	assert := assert.New(t)

	m := make(map[string]int)

	pairs := make([]*P, 1_000_000)

	for i := range pairs {
		pairs[i] = &P{
			k: strconv.Itoa(i),
			v: v2.Int(),
		}
	}

	for _, p := range pairs {
		m[p.k] = p.v
	}

	for _, p := range pairs {
		assert.Equal(p.v, m[p.k])
	}
}
