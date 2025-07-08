package main

import (
	"amoncusir/example/pkg/endpoint"
	"amoncusir/example/pkg/router"
	"amoncusir/example/pkg/service"
	"amoncusir/example/pkg/service/instance"
	"amoncusir/example/pkg/service/scheduler"
	"amoncusir/example/pkg/tools/unres"
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "main: ", log.LstdFlags)
)

func main() {
	uri, err := unres.ExtractFromUri("tcp://0.0.0.0:9090")
	if err != nil {
		logger.Fatal("Invalid URI")
	}

	e, err := endpoint.NewSocket(uri)

	if err != nil {
		logger.Fatal("error when creates socket")
	}

	srv := service.Scheduled(scheduler.RoundRobin())

	rt := router.Broadcast(
		// router.Logger("log: "),
		router.Direct(srv),
	)

	if u, err := unres.ExtractFromUri("tcp://0.0.0.0:9000"); err != nil {
		logger.Fatal("Invalid Instance URI")
	} else {
		srv.AddInstance(instance.New(u))
	}

	e.Listen(rt)

	statusChan := e.Status()

	for status := range statusChan {

		logger.Printf("New endpoint status update: %d", status)

		if status == endpoint.STATUS_CLOSED {
			break
		}
	}
}
