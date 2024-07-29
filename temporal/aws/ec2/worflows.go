package ec2

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func DeployEc2Workflow(ctx workflow.Context, input CreateEc2Input) (CreateEc2Output, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour,
		HeartbeatTimeout:    time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    5 * time.Second,
			BackoffCoefficient: 1.3,
			MaximumInterval:    10 * time.Second,
		},
	})

	var result CreateEc2Output

	// apply and push docker
	if err := workflow.ExecuteActivity(ctx, CreateEc2Activity, input).Get(ctx, &result); err != nil {
		return result, err
	}

	return result, nil
}

func DestroyEc2Workflow(ctx workflow.Context, input DestroyEc2Input) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour,
		HeartbeatTimeout:    time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    10 * time.Second,
			BackoffCoefficient: 1.3,
			MaximumInterval:    2 * time.Minute,
		},
	})

	if err := workflow.ExecuteActivity(ctx, DestroyEc2Activity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}
