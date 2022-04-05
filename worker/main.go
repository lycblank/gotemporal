package main

import (
	"context"
	"github/lycblank/gotemporal/workflow"
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "test", worker.Options{
		BackgroundActivityContext: context.WithValue(context.Background(), "client", c),
	})

	w.RegisterActivity(workflow.UploadProgressActivity)
	w.RegisterWorkflow(workflow.UploadProgressWorkflow)
	w.RegisterWorkflow(workflow.RecvProgressWorkflow)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
