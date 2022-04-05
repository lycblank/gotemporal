package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"go.temporal.io/sdk/client"
	goworkflow "go.temporal.io/sdk/workflow"
	"time"
)

func UploadProgressWorkflow(ctx goworkflow.Context, signalWorkflowID string) error {
	ao := goworkflow.ActivityOptions{
		ScheduleToStartTimeout: 10 * time.Hour,
		StartToCloseTimeout:    5 * time.Hour,
		ScheduleToCloseTimeout: 10 * time.Hour,
		HeartbeatTimeout:       0,
	}
	memo, _ := json.Marshal(goworkflow.GetInfo(ctx).Memo)
	fmt.Println("memo:", string(memo))

	actCtx := goworkflow.WithActivityOptions(ctx, ao)
	goworkflow.ExecuteActivity(actCtx, UploadProgressActivity, signalWorkflowID).Get(ctx, nil)
	return nil
}

func RecvProgressWorkflow(ctx goworkflow.Context) error {
	cwo := goworkflow.ChildWorkflowOptions{
		WorkflowExecutionTimeout: 10 * time.Hour,
		WorkflowTaskTimeout:      time.Hour,
		Memo: map[string]interface{}{
			"workflowID": goworkflow.GetInfo(ctx).WorkflowExecution.ID,
		},
	}

	childCtx := goworkflow.WithChildOptions(ctx, cwo)
	feature := goworkflow.ExecuteChildWorkflow(childCtx, UploadProgressWorkflow, goworkflow.GetInfo(ctx).WorkflowExecution.ID)
	fs := []goworkflow.ChildWorkflowFuture{feature}
	ch := goworkflow.GetSignalChannel(ctx, goworkflow.GetInfo(ctx).WorkflowExecution.ID)
	for len(fs) > 0 {
		fmt.Println("for =====")
		s := goworkflow.NewSelector(ctx)
		for i, f := range fs {
			index := i
			s.AddFuture(f, func(f goworkflow.Future){
				err := f.Get(ctx, nil)
				if err != nil {
					fmt.Println(err)
				}
				fs = append(fs[:index], fs[index+1:]...)
			})
		}
		s.AddReceive(ch, func(c goworkflow.ReceiveChannel, more bool) {
			if c == nil || !more {
				return
			}
			var val int32
			if ok := c.Receive(ctx, &val); ok {
				fmt.Println("progress:", val)
			}
		})
		s.Select(ctx)
	}
	return nil
}

func UploadProgressActivity(ctx context.Context, signalWorkflowID string) error {
	c := ctx.Value("client").(client.Client)
	for i := 0; i < 100; i++ {
		if err := c.SignalWorkflow(ctx, signalWorkflowID, "", signalWorkflowID, i+1); err != nil {
			fmt.Println(err)
		}
		time.Sleep(200 * time.Millisecond)
	}
	return nil
}
