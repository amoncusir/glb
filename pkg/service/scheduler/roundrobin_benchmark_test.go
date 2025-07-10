package scheduler_test

import (
	"amoncusir/example/pkg/service/instance"
	"amoncusir/example/pkg/service/scheduler"
	"amoncusir/example/pkg/tools/unres"
	"fmt"
	"testing"
)

func generateInstances(n int) []instance.Instance {
	ins := make([]instance.Instance, n)

	for i := range ins {
		u, _ := unres.ExtractFromUri(fmt.Sprintf("tcp://%d", i))
		ins[i] = instance.New(u)
	}

	return ins
}

func BenchmarkRoundRobinOneThread(b *testing.B) {
	ins := generateInstances(50)

	b.SetParallelism(1)

	b.Run("Atomic Int counter implementation", func(b *testing.B) {
		sh := scheduler.RoundRobin()
		for b.Loop() {
			sh.Select(ins)
		}
	})

	b.Run("Buffered Channel counter implementation with double of instance size", func(b *testing.B) {
		sh := scheduler.ChanneledRoundRobin(100)
		for b.Loop() {
			sh.Select(ins)
		}
	})

	b.Run("Buffered Channel counter implementation with half of instance size", func(b *testing.B) {
		sh := scheduler.ChanneledRoundRobin(25)
		for b.Loop() {
			sh.Select(ins)
		}
	})
}

func BenchmarkRoundRobinParallel(b *testing.B) {
	ins := generateInstances(50)

	b.Run("Atomic Int counter implementation", func(b *testing.B) {
		sh := scheduler.RoundRobin()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				sh.Select(ins)
			}
		})
	})

	b.Run("Buffered Channel counter implementation with double of instance size", func(b *testing.B) {
		sh := scheduler.ChanneledRoundRobin(100)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				sh.Select(ins)
			}
		})
	})

	b.Run("Buffered Channel counter implementation with half of instance size", func(b *testing.B) {
		sh := scheduler.ChanneledRoundRobin(25)
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				sh.Select(ins)
			}
		})
	})
}
