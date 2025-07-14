package bitmapcheck

var BITMAP_MASK []uint64

func init() {
	BITMAP_MASK = make([]uint64, 64)

	for i := range 64 {
		BITMAP_MASK[i] = 1 << i
	}
}

func NewBoolean(size int) *BinBitmap {
	arraySize := size / 64
	if size % 64 > 0 {
		arraySize++
	}

	return &BinBitmap{
		bitmap: make([]uint64, arraySize),
	}
}

type BinBitmap struct {
	bitmap []uint64
}

func (b *BinBitmap) Size() (n int) {
	return len(b.bitmap) * 64
}

func (b *BinBitmap) MarkPosition(n int, v bool) {
	i := n / 64

	chunk := b.bitmap[i]
	mask := BITMAP_MASK[n%64]

	if v {
		b.bitmap[i] = chunk | mask
	} else {
		b.bitmap[i] = chunk & ^mask
	}
}

func (b *BinBitmap) GetMark(n int) bool {
	i := n / 64

	chunk := b.bitmap[i]
	mask := BITMAP_MASK[n%64]

	return (chunk & mask) > 0
}
