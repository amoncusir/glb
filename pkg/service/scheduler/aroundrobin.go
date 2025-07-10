package scheduler

import (
	"amoncusir/example/pkg/service/instance"
	"runtime"
)

func AsyncRoundRobin(buffer int) Scheduler {
	generated := make(chan int, buffer)
	canceled := make(chan int)

	go func() {
		var index int
		for {
			select {
			case <-canceled:
				return
			default:
				generated <- index
				index++
			}
		}
	}()

	s := &asyncRoundRobin{
		generated: generated,
	}

	// Another more efficient way to do this is add a Close() method to clean up resources, but implies modify the contract
	runtime.AddCleanup(s, func(ch chan int) {
		ch <- 0
	}, canceled)

	return s
}

type asyncRoundRobin struct {
	generated chan int
}

func (s *asyncRoundRobin) Select(inst []instance.Instance) instance.Instance {
	var i int = <-s.generated

	return inst[i%len(inst)]
}
