package scheduler

import (
	"amoncusir/example/pkg/service/instance"
)

type Scheduler interface {
	// Select Must return a slice for the selected instances to reply to a message
	Select(inst []instance.Instance) instance.Instance
}
