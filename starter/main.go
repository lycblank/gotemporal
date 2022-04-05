package main

import (
	"context"
	"fmt"
	"github/lycblank/gotemporal/workflow"
	"go.temporal.io/sdk/client"
	"log"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	// This workflow ID can be user business logic identifier as well.
	workflow1Options := client.StartWorkflowOptions{
		TaskQueue: "test",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflow1Options, workflow.RecvProgressWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow1", err)
	} else {
		log.Println("Started workflow1", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
	}

	we.Get(context.TODO(), nil)
	fmt.Println("start end")
}
