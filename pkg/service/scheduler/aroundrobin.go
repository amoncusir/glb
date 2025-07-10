package scheduler

import (
	"amoncusir/example/pkg/service/instance"
	"runtime"
)

// This RoundRobin implementation using channel performs worst than the other, wich uses Atomic Counters
func ChanneledRoundRobin(buffer int) Scheduler {
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

	s := &channelRoundRobin{
		generated: generated,
	}

	// Another more efficient way to do this is add a Close() method to clean up resources, but implies modify the contract
	runtime.AddCleanup(s, func(ch chan int) {
		ch <- 0
	}, canceled)

	return s
}

type channelRoundRobin struct {
	generated chan int
}

func (s *channelRoundRobin) Select(inst []instance.Instance) instance.Instance {
	if len(inst) <= 0 {
		panic("empty instance slice")
	}

	var i int = <-s.generated

	return inst[i%len(inst)]
}
