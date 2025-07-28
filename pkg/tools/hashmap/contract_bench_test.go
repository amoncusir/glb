package hashmap_test

import (
	"amoncusir/example/pkg/tools/hashmap/openadd"
	"amoncusir/example/pkg/tools/hashmap/swiss"
	"strconv"
	"testing"
)

type entry struct {
	key   string
	value int64
}

const (
	ENTRIES_LENGTH = 1_000_000
)

var entries []*entry

func init() {
	entries = make([]*entry, ENTRIES_LENGTH)

	for i := range ENTRIES_LENGTH {
		entries[i] = &entry{
			key:   strconv.Itoa(i),
			value: int64(i),
		}
	}
}

func BenchmarkStdMap(b *testing.B) {
	b.SetParallelism(1)

	b.Run("Set", func(b *testing.B) {
		m := make(map[string]int64, 256)
		i := 0

		for ; b.Loop(); i++ {
			m[strconv.Itoa(i)] = int64(i)
			i++
		}
	})

	b.Run("Get", func(b *testing.B) {
		m := make(map[string]int64, 256)

		for _, e := range entries {
			m[e.key] = e.value
		}

		i := 0

		for ; b.Loop(); i++ {
			e := entries[i%ENTRIES_LENGTH]
			_ = m[e.key]
		}
	})

	b.Run("Get Empty", func(b *testing.B) {
		m := make(map[string]int64, 256)

		for _, e := range entries {
			m[e.key] = e.value
		}

		i := 0

		for ; b.Loop(); i++ {
			_ = m[strconv.Itoa(-i)]
		}
	})
}

func BenchmarkOpenHashmap(b *testing.B) {
	b.SetParallelism(1)

	b.Run("Set", func(b *testing.B) {
		m := openadd.New(256, 0.7, 2.)
		i := 0

		for ; b.Loop(); i++ {
			m.Set(strconv.Itoa(i), int64(i))
			i++
		}
	})

	b.Run("Get", func(b *testing.B) {
		m := openadd.New(256, 0.7, 2.)

		for _, e := range entries {
			m.Set(e.key, e.value)
		}

		i := 0

		for ; b.Loop(); i++ {
			e := entries[i%ENTRIES_LENGTH]
			_ = m.Get(e.key)
		}
	})

	b.Run("Get Empty", func(b *testing.B) {
		m := openadd.New(256, 0.7, 2.)

		for _, e := range entries {
			m.Set(e.key, e.value)
		}

		i := 0

		for ; b.Loop(); i++ {
			_ = m.Get(strconv.Itoa(-i))
		}
	})
}

func BenchmarkSwissHashmap(b *testing.B) {
	b.SetParallelism(1)

	b.Run("Set", func(b *testing.B) {
		m := swiss.NewM8(256, 0.7, 2.)
		i := 0

		for ; b.Loop(); i++ {
			m.Set(strconv.Itoa(i), int64(i))
			i++
		}
	})

	b.Run("Get", func(b *testing.B) {
		m := swiss.NewM8(256, 0.7, 2.)

		for _, e := range entries {
			m.Set(e.key, e.value)
		}

		i := 0

		for ; b.Loop(); i++ {
			e := entries[i%ENTRIES_LENGTH]
			_ = m.Get(e.key)
		}
	})

	b.Run("Get Empty", func(b *testing.B) {
		m := swiss.NewM8(256, 0.7, 2.)

		for _, e := range entries {
			m.Set(e.key, e.value)
		}

		i := 0

		for ; b.Loop(); i++ {
			_ = m.Get(strconv.Itoa(-i))
		}
	})
}
