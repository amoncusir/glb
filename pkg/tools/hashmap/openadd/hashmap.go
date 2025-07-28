package openadd

import (
	"amoncusir/example/pkg/tools/hashmap"

	"github.com/cespare/xxhash"
)

type HMPair struct {
	key   []byte
	value any
}

func New(initCap uint64, loadFactor, resizeFactor float64) *Hashmap {
	return &Hashmap{
		loadFactor:   loadFactor,
		resizeFactor: resizeFactor,
		items:        0,
		capacity:     initCap,
		content:      make([]*HMPair, initCap),
	}
}

// open-addressed hastable simple implementation
type Hashmap struct {
	loadFactor   float64
	resizeFactor float64
	items        int
	capacity     uint64

	content []*HMPair
}

func hashKey(key []byte, keySize uint64) uint64 {
	return xxhash.Sum64(key) % keySize
}

func (m *Hashmap) Set(key any, v any) {
	// Check if the loadFactor stills lower
	cf := float64(m.items) / float64(m.capacity)
	if m.loadFactor < cf {
		// Resize hashmap!
		m.resize(uint64(float64(m.capacity) * m.resizeFactor))
	}

	m.simpleSet(hashmap.ParseKey(key), v)
}

func (m *Hashmap) resize(cap uint64) {
	pairs := m.Entries()

	m.capacity = cap
	m.content = make([]*HMPair, cap)
	m.items = 0

	for _, p := range pairs {
		m.simpleSet(p.key, p.value)
	}
}

func (m *Hashmap) simpleSet(akey any, v any) {
	key := hashmap.ParseKey(akey)
	i := hashKey(key, m.capacity)

	var p *HMPair
	for ; ; i++ {
		if i >= m.capacity {
			i = 0
		}

		// Check if the value already exists
		p = m.content[i]

		if p == nil {
			break
		}

		if hashmap.Beq(p.key, key) {
			p.value = v
			return
		}
	}

	p = &HMPair{
		key:   key[:],
		value: v,
	}

	m.content[i] = p
	m.items++
}

func (m *Hashmap) Get(akey any) any {
	key := hashmap.ParseKey(akey)
	i := hashKey(key, m.capacity)
	lap := i - 1

	var p *HMPair
	for ; i != lap; i++ {
		if i >= m.capacity {
			i = 0
		}

		// Check if the value already exists
		p = m.content[i]

		if p == nil {
			break
		}

		if hashmap.Beq(p.key, key) {
			return p.value
		}
	}

	return nil
}

func (m *Hashmap) Del(akey any) any {
	key := hashmap.ParseKey(akey)
	i := hashKey(key, m.capacity)
	lap := i - 1

	var p *HMPair
	for ; i != lap; i++ {
		if i >= m.capacity {
			i = 0
		}

		// Check if the value already exists
		p = m.content[i]

		if p == nil {
			break
		}

		if hashmap.Beq(p.key, key) {
			m.content[i] = nil
			m.items--
			return p.value
		}
	}

	return nil
}

func (m *Hashmap) Size() int {
	return m.items
}

func (m *Hashmap) Entries() []*HMPair {
	filledValues := make([]*HMPair, m.items)
	i := 0

	for _, p := range m.content {
		if p != nil {
			filledValues[i] = p
			i++
		}
	}

	return filledValues
}
