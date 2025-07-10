package router

import (
	"amoncusir/example/pkg/router/filter"
	"amoncusir/example/pkg/types"
	"context"
	"errors"
)

func PreFiltered(next Router, filters ...filter.PreFilter) Router {
	return &filterRouter{
		next:    next,
		filters: filters,
	}
}

type filterRouter struct {
	next    Router
	filters []filter.PreFilter
}

func (r *filterRouter) RouteRequest(ctx context.Context, req types.RequestConn) (err error) {

	for _, f := range r.filters {
		if e := f.Accept(ctx, req); e != nil {
			err = errors.Join(err, e)
		}
	}

	if err == nil {
		err = r.next.RouteRequest(ctx, req)
	}

	return err
}
