package filter

import (
	"amoncusir/example/pkg/types"
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

const (
	PRESSURE_SIZE  = 16                 // 2^16 = 65536 values
	LAST_TIME_SIZE = 64 - PRESSURE_SIZE // 2^48 = 281474976710656 values
)

func TBRatelimit(n uint16, inMillis uint) PreFilter {
	return newtbRatelimit(n, inMillis)
}

func newtbRatelimit(capacity uint16, window uint) *tbRatelimit {
	return &tbRatelimit{
		capacity: int64(capacity),
		window:   int64(window),
		weight:   1,
		buckets:  &sync.Map{},
	}
}

// TokenBucket implementation to filter by reate limit
// int64 type prevent casting and type conversion during the compute
type tbRatelimit struct {
	capacity int64
	window   int64
	weight   int64
	buckets  *sync.Map
}

// Accept implements PreFilter.
func (r *tbRatelimit) Accept(ctx context.Context, req types.RequestConn) error {
	now := int64(time.Now().UnixMilli()) & 0x0000_FFFFFFFFFFFF // Set the same precision
	b := r.getBucketByReq(req)

	bucket := b.Load()

	lastPressure := bucket >> LAST_TIME_SIZE
	elapsed := now - (bucket & 0x0000_FFFFFFFFFFFF)

	compression := elapsed / r.window
	pressure := min(r.capacity, lastPressure+compression)

	if pressure <= 0 {
		// Not Allowed also not saved as modifier
		return errors.New("blocked petition by rate limit")
	}

	bucket = ((pressure - r.weight) << LAST_TIME_SIZE) | now

	b.Store(bucket)

	return nil
}

// Uses the first (BE) 16 bits for the counter and last 48 bits for the last acceded time
func (r *tbRatelimit) getBucketByReq(req types.RequestConn) *atomic.Int64 {
	// TODO: The current req.RemoteAddr().String() get the IP + Port of the remote. Must be changed!
	b, _ := r.buckets.LoadOrStore(req.RemoteAddr().String(), &atomic.Int64{})
	return b.(*atomic.Int64)
}
