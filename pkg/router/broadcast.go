package router

import (
	"amoncusir/example/pkg/types"
	"context"
	"errors"
)

func Broadcast(routers ...Router) Router {
	return &broadcastRouter{
		routers: routers,
	}
}

type broadcastRouter struct {
	routers []Router
}

func (r *broadcastRouter) RouteRequest(ctx context.Context, req types.RequestConn) error {
	var err error

	for _, s := range r.routers {
		if e := s.RouteRequest(ctx, req); e != nil {
			err = errors.Join(err, e)
		}
	}

	return err
}
