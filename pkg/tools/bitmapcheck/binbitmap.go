package bitmapcheck

const WORD_SIZE = 32 << (^uint(0) >> 63)

var (
	bitmapMask []uint
)

func init() {
	bitmapMask = make([]uint, WORD_SIZE)

	for i := range WORD_SIZE {
		bitmapMask[i] = 1 << i
	}
}

func NewBoolean(size int) *BinBitmap {
	arraySize := size / WORD_SIZE
	if size%WORD_SIZE > 0 {
		arraySize++
	}

	return &BinBitmap{
		bitmap: make([]uint, arraySize),
	}
}

type BinBitmap struct {
	bitmap []uint
}

func (b *BinBitmap) Size() (n int) {
	return len(b.bitmap) * WORD_SIZE
}

func (b *BinBitmap) MarkPosition(n int, v bool) {
	i := n / WORD_SIZE

	word := b.bitmap[i]
	mask := bitmapMask[n%WORD_SIZE]

	if v {
		b.bitmap[i] = word | mask
	} else {
		b.bitmap[i] = word & ^mask
	}
}

func (b *BinBitmap) GetMark(n int) bool {
	i := n / WORD_SIZE

	word := b.bitmap[i]
	mask := bitmapMask[n%WORD_SIZE]

	return (word & mask) > 0
}
