package deployment

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func CreateDeploymentWorkflow(ctx workflow.Context, input CreateDeploymentInput) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour,
		HeartbeatTimeout:    time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    5 * time.Second,
			BackoffCoefficient: 1.3,
			MaximumInterval:    10 * time.Second,
			MaximumAttempts:    1,
		},
	})

	// Create the Namespace
	if err := workflow.ExecuteActivity(ctx, CreateDeploymentActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

func DestroyDeploymentWorkflow(ctx workflow.Context, input DestroymentInput) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour,
		HeartbeatTimeout:    time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    10 * time.Second,
			BackoffCoefficient: 1.3,
			MaximumInterval:    2 * time.Minute,
			MaximumAttempts:    1,
		},
	})

	if err := workflow.ExecuteActivity(ctx, DestroyDeploymentActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}
