package hashmap

import (
	"amoncusir/example/pkg/tools/hashmap"

	"github.com/cespare/xxhash"
)

func NewM8(cap uint64, loadFactor, resizeFactor float64) *HashmapM8 {
	partitions := cap / PARTITION_LENGTH

	if partitions%PARTITION_LENGTH > 0 {
		partitions++
	}

	return &HashmapM8{
		loadFactor:   loadFactor,
		resizeFactor: resizeFactor,

		items:      0,
		partitions: partitions,
		capacity:   partitions * PARTITION_LENGTH,

		content:  makeContent(partitions),
		metadata: makeMetadata(partitions),
	}
}

const (
	PARTITION_LENGTH        = 8
	LOW_MASK_BITS           = 8                         // Myb 4?
	LOW_MASK         uint64 = 0x00_00_00_00_00_00_00_FF // Myb 0x00_00_00_00_00_00_00_0F?
	HIG_MASK         uint64 = ^LOW_MASK
)

type HMPair struct {
	key   []byte
	value any
}

// open-addressed with aproximation of swiss impl (dont' use SIMD operations and has potential for further optimization)
type HashmapM8 struct {
	loadFactor   float64
	resizeFactor float64

	items      int
	capacity   uint64
	partitions uint64

	content  [][]*HMPair
	metadata [][]uint16
}

func (m *HashmapM8) Set(key any, v any) {
	// Check if the loadFactor stills lower
	cf := float64(m.items) / float64(m.capacity)
	if m.loadFactor < cf {
		// Resize hashmap!
		m.resize(uint64(float64(m.capacity) * m.resizeFactor))
	}

	m.internalSet(hashmap.ParseKey(key), v)
}

func hashKeyHL(key []byte) (hk uint64, lk uint8) {
	hk = xxhash.Sum64(key)

	// hk = ((hash & HIG_MASK) >> LOW_MASK_BITS)
	// lk = uint8(hk & LOW_MASK)
	lk = uint8(hk % 255)
	return
}

func (m *HashmapM8) resize(cap uint64) {
	pairs := m.Entries()

	partitions := cap / PARTITION_LENGTH

	if partitions%PARTITION_LENGTH > 0 {
		partitions++
	}

	m.items = 0
	m.partitions = partitions
	m.capacity = partitions * PARTITION_LENGTH
	m.content = makeContent(partitions)
	m.metadata = makeMetadata(partitions)

	for _, p := range pairs {
		m.internalSet(p.key, p.value)
	}
}

func makeContent(p uint64) (content [][]*HMPair) {
	content = make([][]*HMPair, p)

	for i := range content {
		content[i] = make([]*HMPair, PARTITION_LENGTH)
	}
	return
}

func makeMetadata(p uint64) (metadata [][]uint16) {
	metadata = make([][]uint16, p)

	for i := range metadata {
		metadata[i] = make([]uint16, PARTITION_LENGTH)
	}
	return
}

// First, search if exists
// If exist, replace it
// If not, find the first empty space, place it and add 1 to items count
func (m *HashmapM8) internalSet(key []byte, v any) {
	// Extract keys and module to the right size
	hk, lk := hashKeyHL(key)
	mainIndex := hk % m.partitions
	chukIndex := lk % PARTITION_LENGTH

	// Happy path
	container := m.content[mainIndex]
	pair := container[chukIndex]

	// Iterate over the slotes to try to find an empty one or if exist, iterates the existent pairs
	for pair != nil {

		// First, check if the current pair is the same
		if hashmap.Beq(pair.key, key) {
			pair.value = v
			return
		}

		// If not, find the next
		chukIndex++

		// If the chunk no has more items, reset to 0 and try the next one
		if chukIndex >= PARTITION_LENGTH {
			chukIndex = 0
			mainIndex++

			// Also if we reach the end of the partitions, go to the first one
			if mainIndex >= m.partitions {
				mainIndex = 0
			}

			container = m.content[mainIndex]
		}

		pair = container[chukIndex]
	}

	pair = &HMPair{
		key:   key[:],
		value: v,
	}

	m.metadata[mainIndex][chukIndex] = 0xFFFF & uint16(lk)
	container[chukIndex] = pair
	m.items++
}

func (m *HashmapM8) Get(akey any) any {
	// Extract keys and module to the right size
	key := hashmap.ParseKey(akey)
	hk, lk := hashKeyHL(key)
	mainIndex := hk % m.partitions
	chukIndex := lk % PARTITION_LENGTH

	// Happy path
	container := m.content[mainIndex]
	pair := container[chukIndex]

	// Iterate over the slotes to try to find an empty one or if exist, iterates the existent pairs
	for pair != nil {

		// First, check if the current pair is the same
		if hashmap.Beq(pair.key, key) {
			return pair.value
		}

		// If not, find the next
		chukIndex++

		// If the chunk no has more items, reset to 0 and try the next one
		if chukIndex >= PARTITION_LENGTH {
			chukIndex = 0
			mainIndex++

			// Also if we reach the end of the partitions, go to the first one
			if mainIndex >= m.partitions {
				mainIndex = 0
			}

			container = m.content[mainIndex]
		}

		// iterate over metadata values until next are zero or matches the low key
		for meta := m.metadata[mainIndex][chukIndex]; meta != 0 && uint8(meta) != lk; {
			// If not, find the next
			chukIndex++

			// If the chunk no has more items, reset to 0 and try the next one
			if chukIndex >= PARTITION_LENGTH {
				chukIndex = 0
				mainIndex++

				// Also if we reach the end of the partitions, go to the first one
				if mainIndex >= m.partitions {
					mainIndex = 0
				}

				container = m.content[mainIndex]
			}

			meta = m.metadata[mainIndex][chukIndex]
		}

		pair = container[chukIndex]
	}

	return nil
}

func (m *HashmapM8) Del(akey any) any {
	// Extract keys and module to the right size
	key := hashmap.ParseKey(akey)
	hk, lk := hashKeyHL(key)
	mainIndex := hk % m.partitions
	chukIndex := lk % PARTITION_LENGTH

	// Happy path
	container := m.content[mainIndex]
	pair := container[chukIndex]

	// Iterate over the slotes to try to find an empty one or if exist, iterates the existent pairs
	for pair != nil {

		// First, check if the current pair is the same
		if hashmap.Beq(pair.key, key) {
			container[chukIndex] = nil
			m.items--
			return pair.value
		}

		// If not, find the next
		chukIndex++

		// If the chunk no has more items, reset to 0 and try the next one
		if chukIndex >= PARTITION_LENGTH {
			chukIndex = 0
			mainIndex++

			// Also if we reach the end of the partitions, go to the first one
			if mainIndex >= m.partitions {
				mainIndex = 0
			}

			container = m.content[mainIndex]
		}

		// iterate over metadata values until next are zero or matches the low key
		for meta := m.metadata[mainIndex][chukIndex]; meta != 0 && uint8(meta) != lk; {
			// If not, find the next
			chukIndex++

			// If the chunk no has more items, reset to 0 and try the next one
			if chukIndex >= PARTITION_LENGTH {
				chukIndex = 0
				mainIndex++

				// Also if we reach the end of the partitions, go to the first one
				if mainIndex >= m.partitions {
					mainIndex = 0
				}

				container = m.content[mainIndex]
			}

			meta = m.metadata[mainIndex][chukIndex]
		}

		pair = container[chukIndex]
	}

	return nil
}

func (m *HashmapM8) Size() int {
	return m.items
}

func (m *HashmapM8) Entries() []*HMPair {
	filledValues := make([]*HMPair, m.items)
	v := 0

	for i := range m.content {
		for ii := range PARTITION_LENGTH {
			if p := m.content[i][ii]; p != nil {
				filledValues[v] = p
				v++
			}
		}
	}

	return filledValues
}
