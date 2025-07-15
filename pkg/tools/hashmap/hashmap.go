package hashmap

import (
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

func beq(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (m *Hashmap) Set(key []byte, v any) {
	// Check if the loadFactor stills lower
	cf := float64(m.items) / float64(m.capacity)
	if m.loadFactor < cf {
		// Resize hashmap!
		m.resize(uint64(float64(m.capacity) * m.resizeFactor))
	}

	m.simpleSet(key, v)
}

func (m *Hashmap) resize(cap uint64) {
	pairs := m.Pairs()
	new := &Hashmap{
		loadFactor:   m.loadFactor,
		resizeFactor: m.resizeFactor,
		items:        0,
		capacity:     cap,
		content:      make([]*HMPair, cap),
	}

	for _, p := range pairs {
		new.simpleSet(p.key, p.value)
	}

	m.capacity = cap
	m.content = new.content
}

func (m *Hashmap) simpleSet(key []byte, v any) {
	i := hashKey(key, m.capacity)

	var p *HMPair
	for ; ; i++ {
		if i > m.capacity {
			i = 0
		}

		// Check if the value already exists
		p = m.content[i]

		if p == nil {
			break
		}

		if beq(p.key, key) {
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

func (m *Hashmap) Get(key []byte) any {
	i := hashKey(key, m.capacity)
	lap := i - 1

	var p *HMPair
	for ; i != lap; i++ {
		if i > m.capacity {
			i = 0
		}

		// Check if the value already exists
		p = m.content[i]

		if p == nil {
			break
		}

		if beq(p.key, key) {
			return p.value
		}
	}

	return nil
}

func (m *Hashmap) Del(key []byte) any {
	i := hashKey(key, m.capacity)
	lap := i - 1

	var p *HMPair
	for ; i != lap; i++ {
		if i > m.capacity {
			i = 0
		}

		// Check if the value already exists
		p = m.content[i]

		if p == nil {
			break
		}

		if beq(p.key, key) {
			m.content[i] = nil
			return p.value
		}
	}

	return nil
}

func (m *Hashmap) Size() int {
	return m.items
}

func (m *Hashmap) Pairs() []*HMPair {
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
