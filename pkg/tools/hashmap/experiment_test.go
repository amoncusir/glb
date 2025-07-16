package hashmap

import (
	v2 "math/rand/v2"
	"slices"
	"testing"

	"github.com/cespare/xxhash"
)

func generateRandomKey() uint64 {
	len := v2.Int64() % 64
	key := make([]byte, len)

	for i := range len {
		key[i] = byte(v2.Int32() % 255)
	}

	return xxhash.Sum64(key)
}

func TestHashBitColision(t *testing.T) {

	splits := 8
	wordsize := 8
	tests := 32
	iterations := 2_000_000
	colisions := 0
	falseNegative := 0
	zeroHash := 0

	gLowk := func(k uint64) uint64 {
		// return ((k * 0x9dc5) & 0x00_00_00_00_00_00_00_FF)
		return (k % 255) // & 0x00_00_00_00_00_00_00_FF
		// return k & 0x0000_0000_0000_FFFF
		// return k & 0x000_000_000_000_FFF
	}

	matchLowK := func(sum, lowk uint64) bool {
		for i := range splits {
			plowk := lowk << (i * wordsize)
			match := sum & plowk
			if match == plowk {
				return true
			}
		}

		return false
	}

	for range iterations {
		sum := uint64(0)
		keys := make([]uint64, splits)

		for i := range splits {
			k := generateRandomKey()
			lowk := gLowk(k)

			keys[i] = k

			sum = sum | (lowk << (i * wordsize))
		}

		for _, v := range keys {
			lowk := gLowk(v)

			if !matchLowK(sum, lowk) {
				falseNegative++
			}
		}

		for range tests {
			k := generateRandomKey()
			lowk := gLowk(k)

			if lowk == 0 {
				zeroHash++
			}

			if matchLowK(sum, lowk) {
				if slices.Index(keys, k) < 0 {
					colisions++
				}
			}
		}
	}

	t.Logf("Total iterations %d with %d colisions and false negatives %d", iterations*tests, colisions, falseNegative)
	t.Logf("Zero hashes %d", zeroHash)
	t.Logf("Current precision: %.2f", 1.-float64(colisions)/float64(iterations*tests))
}
