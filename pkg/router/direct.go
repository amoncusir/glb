package router

import (
	"amoncusir/example/pkg/service"
	"amoncusir/example/pkg/types"
	"context"
)

func Direct(service service.Service) Router {
	return &directRouter{
		service: service,
	}
}

type directRouter struct {
	service service.Service
}

func (r *directRouter) RouteRequest(ctx context.Context, req types.RequestConn) error {
	return r.service.Reply(ctx, req)
}
