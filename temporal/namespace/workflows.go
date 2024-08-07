package namespace

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func CreateNamespaceWorkflow(ctx workflow.Context, input CreateNamespaceInput) error {
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
	if err := workflow.ExecuteActivity(ctx, CreateNamespaceActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

func DestroyNamespaceWorkflow(ctx workflow.Context, input DestroyNamespaceInput) error {
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

	if err := workflow.ExecuteActivity(ctx, DestroyNamespaceActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}
