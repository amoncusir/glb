package router

import (
	"amoncusir/example/pkg/types"
	"context"
	"log"
	"os"
)

func Logger(name string) Router {
	return &logRouter{
		log: log.New(os.Stdout, name, log.LstdFlags),
	}
}

type logRouter struct {
	log *log.Logger
}

func (r *logRouter) RouteRequest(ctx context.Context, req types.RequestConn) error {
	r.log.Printf("New Connection from: %s", req.RemoteAddr())
	return nil
}
