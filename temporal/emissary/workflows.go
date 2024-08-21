package emissary

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func CreateEmissaryMappingWorkflow(ctx workflow.Context, input CreateEmissaryMappingInput) error {
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
	if err := workflow.ExecuteActivity(ctx, CreateEmissaryMappingActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

func DestroyEmissaryMappingWorkflow(ctx workflow.Context, input DestroyEmissaryMappingInput) error {
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

	if err := workflow.ExecuteActivity(ctx, DestroyEmissaryMappingActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}
