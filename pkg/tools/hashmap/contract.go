package hashmap

type Pair struct {
	key   any
	value any
}

type Hasmap interface {
	Set(key any, v any)
	Get(key any) any
	Del(key any) any

	Size() int
	Entries() []*Pair
}
