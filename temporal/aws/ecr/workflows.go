package ecr

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func ApplyDockerBuildAndPushEcrWorkflow(ctx workflow.Context, input ApplyDockerBuildAndPushEcrInput) (ApplyDockerBuildAndPushEcrOutput, error) {
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

	var result ApplyDockerBuildAndPushEcrOutput

	// apply and push docker
	if err := workflow.ExecuteActivity(ctx, ApplyDockerBuildAndPushEcrActivity, input).Get(ctx, &result); err != nil {
		return result, err
	}

	return result, nil
}

func DestroyDockerBuildAndPushEcrWorkflow(ctx workflow.Context, input DestroyDockerBuildAndPushEcrInput) error {
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

	if err := workflow.ExecuteActivity(ctx, DestroyDockerBuildAndPushEcrActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}
