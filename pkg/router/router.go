package router

import (
	"amoncusir/example/pkg/types"
	"context"
)

// Must be thread safe!
type Router interface {
	RouteRequest(ctx context.Context, req types.RequestConn) error
}
