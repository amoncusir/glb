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

func TBRatelimit(size uint16, recoveryTime, recoveryUnit uint) PreFilter {
	rl := newtbRatelimit(size, recoveryTime)
	rl.recovery = int64(recoveryUnit)
	return rl
}

func newtbRatelimit(capacity uint16, window uint) *tbRatelimit {
	return &tbRatelimit{
		timeRef:   int64(time.Now().UnixMilli()),
		capacity:  int64(capacity),
		window:    int64(window),
		discharge: 1,
		recovery:  1,
		buckets:   &sync.Map{},
	}
}

// TokenBucket implementation to filter by reate limit
// int64 type prevent casting and type conversion during the compute
type tbRatelimit struct {
	timeRef   int64 // Used to reduce the UNIX time reference and give more bits to calculate the elapsed time
	capacity  int64
	window    int64
	recovery  int64
	discharge int64
	buckets   *sync.Map
}

// Accept implements PreFilter.
func (r *tbRatelimit) Accept(ctx context.Context, req types.RequestConn) error {
	now := int64(time.Now().UnixMilli()-r.timeRef) & 0x0000_FFFFFFFFFFFF // Set the same precision bits
	b := r.getBucketByReq(req)

	bucket := b.Load()

	lastPressure := bucket >> LAST_TIME_SIZE
	elapsed := now - (bucket & 0x0000_FFFFFFFFFFFF)

	compression := (elapsed / r.window) * r.recovery
	pressure := min(r.capacity, lastPressure+compression)

	if pressure <= 0 {
		// Not Allowed also not saved as modifier
		return errors.New("blocked petition by rate limit")
	}

	bucket = ((pressure - r.discharge) << LAST_TIME_SIZE) | now

	b.Store(bucket)

	return nil
}

// Uses the first (BE) 16 bits for the counter and last 48 bits for the last acceded time
func (r *tbRatelimit) getBucketByReq(req types.RequestConn) *atomic.Int64 {
	key := req.RemoteIpString()
	b, ok := r.buckets.Load(key)

	if !ok {
		b = &atomic.Int64{}
		b.(*atomic.Int64).Store(r.capacity << LAST_TIME_SIZE & -0x0001_000000000000) // 0xFFFF_000000000000 as uint64 in complement a2
		r.buckets.Store(key, b)
	}

	return b.(*atomic.Int64)
}
