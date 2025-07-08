package scheduler

import (
	"amoncusir/example/pkg/service/instance"
)

type Scheduler interface {
	// Must return an slice for the selected instances to reply the message
	Select(inst []instance.Instance) instance.Instance
}
