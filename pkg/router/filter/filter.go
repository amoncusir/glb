package filter

import (
	"amoncusir/example/pkg/types"
	"context"
)

type PreFilter interface {
	Accept(ctx context.Context, req types.RequestConn) error
}
