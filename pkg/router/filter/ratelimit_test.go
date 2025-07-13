package filter

import (
	"amoncusir/example/test/mock"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRateLimitDenyAfterEmpty(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	req := mock.NewMockRequestConn(ctrl)
	ctx := context.Background()

	req.EXPECT().RemoteIp().Return("0.0.0.0").AnyTimes()

	rl := newtbRatelimit(10, 1_000)

	for range 10 {
		err := rl.Accept(ctx, req)
		assert.Empty(err)
	}

	err := rl.Accept(ctx, req)
	assert.NotEmpty(err)
}

func TestRateLimitDenyedAndAcceptBeforeTime(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	req := mock.NewMockRequestConn(ctrl)
	ctx := context.Background()

	req.EXPECT().RemoteIp().Return("0.0.0.0").AnyTimes()

	rl := newtbRatelimit(10, 1_000)

	for range 10 {
		err := rl.Accept(ctx, req)
		assert.Empty(err)
	}

	assert.NotEmpty(rl.Accept(ctx, req))

	time.Sleep(1 * time.Second)

	assert.Empty(rl.Accept(ctx, req))
	assert.NotEmpty(rl.Accept(ctx, req))
}
