package instance

import (
	"amoncusir/example/pkg/types"
	"context"
	"log"
)

type Status int

const (
	STATUS_HEALTHY   Status = 1
	STATUS_UNHEALTHY Status = -1
	STATUS_UNKNOWN   Status = 0
)

var (
	logger = log.Default()
)

// Instance
// Does not share Instance between Service.
// Each Service must contain only one reference of Instance
type Instance interface {
	Address() string
	Protocol() string
	Status() Status

	AddHealthyCallback(fn func(self Instance)) error
	AddUnhealthyCallback(fn func(self Instance)) error

	Reply(ctx context.Context, req types.RequestConn) error

	// WaitForReplies Wait until all replies are finished and blocks new ones
	WaitForReplies() error

	Watch() error
	Unwatch() error
}

// noCopy may be added to structs which must not be copied
// after the first use.
//
// See https://golang.org/issues/8005#issuecomment-190753527
// for details.
//
// Note that it must not be embedded, due to the Lock and Unlock methods.
type noCopy struct{}
