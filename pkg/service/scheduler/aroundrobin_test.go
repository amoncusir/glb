package scheduler

import (
	"amoncusir/example/pkg/service/instance"
	"amoncusir/example/test/mock"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestARoundRobinFixedSlice(t *testing.T) {
	ctrl := gomock.NewController(t)
	assert := assert.New(t)

	sh := AsyncRoundRobin(5)

	ins := make([]instance.Instance, 10)

	for i := range ins {
		in := mock.NewMockInstance(ctrl)

		in.
			EXPECT().
			Address().
			Return(fmt.Sprintf("%d", i)).
			Times(2)

		ins[i] = in
	}

	for i := range ins {
		in := sh.Select(ins)

		assert.Equal(fmt.Sprintf("%d", i), in.Address())
	}

	for i := range ins {
		in := sh.Select(ins)

		assert.Equal(fmt.Sprintf("%d", i), in.Address())
	}
}

func TestARoundRobinDynamicSlice(t *testing.T) {
	ctrl := gomock.NewController(t)
	assert := assert.New(t)

	sh := AsyncRoundRobin(5)
	ins := make([]instance.Instance, 10)

	for i := range ins {
		in := mock.NewMockInstance(ctrl)

		in.
			EXPECT().
			Address().
			Return(fmt.Sprintf("%d", i)).
			AnyTimes()

		ins[i] = in
	}

	// First, full slice lenght nad move some positions
	assert.Equal("0", sh.Select(ins).Address())
	assert.Equal("1", sh.Select(ins).Address())
	assert.Equal("2", sh.Select(ins).Address())
	assert.Equal("3", sh.Select(ins).Address())

	// Reduce to 4 instances and expect the first one
	assert.Equal("0", sh.Select(ins[:4]).Address())

	// Reduce to 4 instances and expect the second one in different slice
	assert.Equal("5", sh.Select(ins[4:8]).Address())

	// Reduce to 2 instances
	assert.Equal("0", sh.Select(ins[0:2]).Address())
	assert.Equal("1", sh.Select(ins[0:2]).Address())
	assert.Equal("0", sh.Select(ins[0:2]).Address())

	// Full length, last instance
	assert.Equal("9", sh.Select(ins[:]).Address())
}

func TestARoundRobinParallel(t *testing.T) {
	ctrl := gomock.NewController(t)

	sh := AsyncRoundRobin(5)
	ins := make([]instance.Instance, 10)
	n := 4

	for i := range ins {
		in := mock.NewMockInstance(ctrl)

		in.
			EXPECT().
			Address().
			Times(n)

		ins[i] = in
	}

	wg := sync.WaitGroup{}
	wg.Add(n)

	for range n {
		go func() {
			defer wg.Done()
			for range 10 {
				sh.Select(ins).Address()
			}
		}()

	}

	wg.Wait()
}
