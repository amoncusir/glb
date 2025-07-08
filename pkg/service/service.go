package service

import (
	"amoncusir/example/pkg/service/instance"
	"amoncusir/example/pkg/types"
	"context"
	"io"
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "service: ", log.LstdFlags)
)

type Service interface {
	io.Closer

	Reply(ctx context.Context, req types.RequestConn) error

	Instances() []instance.Instance
	AddInstance(inst instance.Instance) error
	DelInstance(inst instance.Instance) error
}

// noCopy may be added to structs which must not be copied
// after the first use.
//
// See https://golang.org/issues/8005#issuecomment-190753527
// for details.
//
// Note that it must not be embedded, due to the Lock and Unlock methods.
type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
