package scheduler

import (
	"amoncusir/example/pkg/service/instance"
	"sync/atomic"
)

func RoundRobin() Scheduler {
	return &roundRobin{
		index: &atomic.Int64{},
	}
}

type roundRobin struct {
	index *atomic.Int64
}

func (s *roundRobin) Select(inst []instance.Instance) instance.Instance {
	return s.nextInstance(inst)
}

func (s *roundRobin) nextInstance(inst []instance.Instance) instance.Instance {
	if len(inst) <= 0 {
		panic("empty instance slice")
	}

	i := s.index.Add(1)
	return inst[(i-1)%int64(len(inst))]
}
