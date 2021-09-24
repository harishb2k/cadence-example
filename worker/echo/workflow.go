package echo

import (
	"context"
	"go.uber.org/cadence/workflow"
	"time"
)

type In struct {
	Name string
}

type Out struct {
	Name   string
	Result string
}

func Activity(ctx context.Context, in *In) (*Out, error) {
	return &Out{
		Name:   in.Name,
		Result: "Success - " + time.Now().String(),
	}, nil
}

func Workflow(ctx workflow.Context, in *In) (*Out, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: 10 * time.Hour,
		StartToCloseTimeout:    10 * time.Hour,
		ScheduleToStartTimeout: 10 * time.Hour,
	})

	var out Out
	err := workflow.ExecuteActivity(ctx, Activity, in).Get(ctx, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
