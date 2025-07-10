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

func TBRatelimit() PreFilter {
	return &rtbRatelimit{}
}

// TokenBucket implementation to filter by reate limit
// uint64 type prevent casting and type conversion during the compute
type rtbRatelimit struct {
	capacity        uint64
	compressionRate uint64 // in milliseconds
	buckets         *sync.Map
}

// Accept implements PreFilter.
func (r *rtbRatelimit) Accept(ctx context.Context, req types.RequestConn) error {
	b := r.getDataByReq(req)
	now := uint64(time.Now().UnixMilli()) & 0x0000_FFFFFFFFFFFF // Set the same precision

	bucket := b.Load()

	lastPressure := bucket >> LAST_TIME_SIZE
	elapsed := now - (bucket & 0x0000_FFFFFFFFFFFF)
	compression := elapsed * r.compressionRate

	pressure := min(r.capacity, lastPressure+compression)

	if pressure <= 0 {
		// Not Allowed also not saved as modifier
		return errors.New("blocked petition by rate limit")
	}

	bucket = ((pressure - 1) << LAST_TIME_SIZE) & now

	b.Store(bucket)

	return nil
}

// Uses the first (BE) 16 bits for the counter and last 48 bits for the last acceded time
func (r *rtbRatelimit) getDataByReq(req types.RequestConn) *atomic.Uint64 {
	// TODO: The current req.RemoteAddr().String() get the IP + Port of the remote. Must be changed!
	b, _ := r.buckets.LoadOrStore(req.RemoteAddr().String(), &atomic.Uint64{})
	return b.(*atomic.Uint64)
}
