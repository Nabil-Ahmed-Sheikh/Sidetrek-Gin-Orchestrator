package dagster

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func CreateDagsterClusterWorkflow(ctx workflow.Context, input CreateDagsterClusterInput) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour,
		HeartbeatTimeout:    time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    5 * time.Second,
			BackoffCoefficient: 1.3,
			MaximumInterval:    10 * time.Second,
			MaximumAttempts:    10,
		},
	})

	// Create the Cluster
	if err := workflow.ExecuteActivity(ctx, CreateDagsterClusterActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

func DestroyDagsterClusterWorkflow(ctx workflow.Context, input DestroyDagsterClusterInput) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour,
		HeartbeatTimeout:    time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    10 * time.Second,
			BackoffCoefficient: 1.3,
			MaximumInterval:    2 * time.Minute,
			MaximumAttempts:    10,
		},
	})

	if err := workflow.ExecuteActivity(ctx, DestroyDagsterClusterActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}
