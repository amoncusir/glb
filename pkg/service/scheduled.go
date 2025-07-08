package service

import (
	"amoncusir/example/pkg/service/instance"
	"amoncusir/example/pkg/service/scheduler"
	"amoncusir/example/pkg/types"
	"context"
	"errors"
	"slices"
	"sync"
)

func Scheduled(alghtm scheduler.Scheduler) Service {

	return &scheduledService{
		alghtm:           alghtm,
		instances:        []instance.Instance{},
		healthyInstances: []instance.Instance{},
		updaterLck:       &sync.RWMutex{},
	}
}

type scheduledService struct {
	_ noCopy

	alghtm scheduler.Scheduler

	instances        []instance.Instance
	healthyInstances []instance.Instance
	updaterLck       *sync.RWMutex
}

// Close implements Service.
func (s *scheduledService) Close() error {
	var pErr error

	for _, ins := range s.instances {
		if err := ins.Close(); err != nil {
			pErr = errors.Join(pErr, err)
		}
	}

	return pErr
}

// Instances implements Service.
func (s *scheduledService) Instances() []instance.Instance {
	return append([]instance.Instance{}, s.instances...)
}

// Instances implements Service.
func (s *scheduledService) AddInstance(inst instance.Instance) error {
	logger.Printf("New instance added: %s\n", inst)

	s.updaterLck.Lock()
	defer s.updaterLck.Unlock()

	if err := inst.AddHealthyCallback(func(self instance.Instance) {
		s.markHealty(self)
	}); err != nil {
		return err
	}

	if err := inst.AddUnhealthyCallback(func(self instance.Instance) {
		s.markUnhealty(self)
	}); err != nil {
		return err
	}

	s.instances = append(s.instances, inst)
	s.healthyInstances = append(s.healthyInstances, inst)

	go inst.Watch()

	return nil
}

func (s *scheduledService) DelInstance(inst instance.Instance) error {
	s.updaterLck.Lock()
	defer s.updaterLck.Unlock()

	indx := slices.IndexFunc(s.instances, func(finst instance.Instance) bool {
		return &finst == &inst
	})

	if indx <= 0 {
		return errors.New("not found any instance")
	}

	s.instances = slices.Delete(s.instances, indx, indx+1)

	if i := slices.IndexFunc(s.healthyInstances, func(finst instance.Instance) bool {
		return &finst == &inst
	}); i >= 0 {
		s.healthyInstances = slices.Delete(s.healthyInstances, i, i+1)
	}

	defer inst.Close()
	defer inst.Unwatch()

	return nil
}

func (s *scheduledService) markHealty(inst instance.Instance) {
	logger.Printf("Mark as healthy instance %v", inst)

	if i := slices.IndexFunc(s.instances, func(finst instance.Instance) bool {
		return finst == inst
	}); i >= 0 {
		s.updaterLck.Lock()
		s.healthyInstances = append(s.healthyInstances, inst)
		s.updaterLck.Unlock()
	}
}

func (s *scheduledService) markUnhealty(inst instance.Instance) {
	logger.Printf("Mark as unhealthy instance %v", inst)

	if i := slices.IndexFunc(s.healthyInstances, func(finst instance.Instance) bool {
		return finst == inst
	}); i >= 0 {
		s.updaterLck.Lock()
		s.healthyInstances = slices.Delete(s.healthyInstances, i, i+1)
		s.updaterLck.Unlock()
	}
}

// Reply implements Service.
func (s *scheduledService) Reply(ctx context.Context, req types.RequestConn) error {
	logger.Printf("Reply msg with total of %d instances and %v healthy", len(s.instances), len(s.healthyInstances))

	s.updaterLck.RLock()
	defer s.updaterLck.RUnlock()

	if len(s.instances) <= 0 {
		return errors.New("no contains any instance")
	}

	if len(s.healthyInstances) <= 0 {
		return errors.New("no contains any healty instance")
	}

	inst := s.alghtm.Select(s.healthyInstances)
	return inst.Reply(ctx, req)
}
