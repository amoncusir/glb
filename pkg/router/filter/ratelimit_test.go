package filter

import (
	"amoncusir/example/test/mock"
	"context"
	"sync"
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

func TestRateLimitReccoverAllCapacity(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	req := mock.NewMockRequestConn(ctrl)
	ctx := context.Background()

	req.EXPECT().RemoteIp().Return("0.0.0.0").AnyTimes()

	rl := newtbRatelimit(10, 1_000/10)

	for range 10 {
		err := rl.Accept(ctx, req)
		assert.Empty(err)
	}

	assert.NotEmpty(rl.Accept(ctx, req))

	time.Sleep(1 * time.Second)

	for range 10 {
		err := rl.Accept(ctx, req)
		assert.Empty(err)
	}

	assert.NotEmpty(rl.Accept(ctx, req))
}

func TestRateLimitOnConcurrentCallsNoDrainCapacity(t *testing.T) {
	iterations := 10
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	req := mock.NewMockRequestConn(ctrl)
	ctx := context.Background()

	req.EXPECT().RemoteIp().Return("0.0.0.0").AnyTimes()

	rl := newtbRatelimit(10, 100/10) // 100 accepts every second

	wg := &sync.WaitGroup{}

	callFn := func() {
		defer wg.Done()

		for range 10 {
			err := rl.Accept(ctx, req)
			assert.Empty(err)
			time.Sleep(100 * time.Millisecond)
		}
	}

	wg.Add(iterations) // Perform 100 calls in one second
	for range iterations {
		go callFn()
	}

	wg.Wait()
}

func TestRateLimitOnConcurrentCallsDrainCapacity(t *testing.T) {
	// Test Parameters
	parallelization := 10_000
	windowTime := 10
	initialBudget := 10
	petitions := 100
	acceptableDeviation := 0.995

	// Arrange
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	req := mock.NewMockRequestConn(ctrl)
	ctx := context.Background()
	errs := make(chan error)
	wg := &sync.WaitGroup{}
	collectedErrors := []error{}

	req.EXPECT().RemoteIp().Return("0.0.0.0").AnyTimes()

	rl := newtbRatelimit(uint16(initialBudget), uint(windowTime/initialBudget)) // 100 accepts every second

	callFn := func() {
		defer wg.Done()
		for range petitions {
			errs <- rl.Accept(ctx, req)
			time.Sleep(time.Duration(windowTime) * time.Millisecond)
		}
	}

	// Act
	wg.Add(parallelization) // Perform 100 calls in one second

	go func() {
		wg.Wait()
		close(errs)
	}()

	for range parallelization {
		go callFn()
	}

	for err := range errs {
		if err != nil {
			collectedErrors = append(collectedErrors, err)
		}
	}

	// Assert
	// This algoritm has eventual consistence and more than one request could be accepted.
	// An acceptable deviation is less then
	totalCaluledErrors := parallelization*petitions - (petitions * initialBudget)
	deviation := 1.0 - acceptableDeviation
	errorsDeviation := int(float64(totalCaluledErrors) * deviation)

	t.Logf("Total errors: %d over calculed: %d and deviation of %d", len(collectedErrors), totalCaluledErrors, errorsDeviation)

	assert.LessOrEqual(len(collectedErrors), totalCaluledErrors+errorsDeviation)
	assert.GreaterOrEqual(len(collectedErrors), totalCaluledErrors-errorsDeviation)
}
