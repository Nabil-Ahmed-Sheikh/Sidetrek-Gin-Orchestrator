package network

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func CreateVpcWorkflow(ctx workflow.Context, input CreateVpcInput) (CreateVpcOutput, error) {
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

	// Create the VPC
	var vpcOutput CreateVpcOutput
	if err := workflow.ExecuteActivity(ctx, CreateVpcActivity, input).Get(ctx, &vpcOutput); err != nil {
		return CreateVpcOutput{}, err
	}

	return CreateVpcOutput{
		VpcID: vpcOutput.VpcID,
	}, nil
}

func DestroyVpcWorkflow(ctx workflow.Context, input DestroyVpcInput) error {
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

	if err := workflow.ExecuteActivity(ctx, DestroyVpcActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

func CreateSubnetsWorkflow(ctx workflow.Context, input CreateSubnetsInput) (CreateSubnetsOutput, error) {
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

	// Create subnets
	var subnetOutput CreateSubnetsOutput
	if err := workflow.ExecuteActivity(ctx, CreateSubnetsActivity, input).Get(ctx, &subnetOutput); err != nil {
		return CreateSubnetsOutput{}, err
	}

	return CreateSubnetsOutput{}, nil
}

func DestroySubnetsWorkflow(ctx workflow.Context, input DestroySubnetsInput) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour,
		HeartbeatTimeout:    time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    10 * time.Second,
			BackoffCoefficient: 1.3,
			MaximumInterval:    2 * time.Minute,
		},
	})

	if err := workflow.ExecuteActivity(ctx, DestroySubnetsActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

func CreateSecurityGroupWorkflow(ctx workflow.Context, input CreateSecurityGroupInput) error {
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

	// Create the SecurityGroup
	if err := workflow.ExecuteActivity(ctx, CreateSecurityGroupActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

func DestroySecurityGroupWorkflow(ctx workflow.Context, input DestroySecurityGroupInput) error {
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

	// Destroy the SecurityGroup
	if err := workflow.ExecuteActivity(ctx, DestroySecurityGroupActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

func CreateElasticIpWorkflow(ctx workflow.Context, input CreateElasticIpInput) error {
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

	// Create the ElasticIp
	if err := workflow.ExecuteActivity(ctx, CreateElasticIpActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

func DestroyElasticIpWorkflow(ctx workflow.Context, input DestroyElasticIpInput) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour,
		HeartbeatTimeout:    time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    10 * time.Second,
			BackoffCoefficient: 1.3,
			MaximumInterval:    2 * time.Minute,
		},
	})

	if err := workflow.ExecuteActivity(ctx, DestroyElasticIpActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

func CreateInternetGatewayWorkflow(ctx workflow.Context, input CreateInternetGatewayInput) error {
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

	// Create the InternetGateway
	if err := workflow.ExecuteActivity(ctx, CreateInternetGatewayActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

func DestroyInternetGatewayWorkflow(ctx workflow.Context, input DestroyInternetGatewayInput) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour,
		HeartbeatTimeout:    time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    10 * time.Second,
			BackoffCoefficient: 1.3,
			MaximumInterval:    2 * time.Minute,
		},
	})

	if err := workflow.ExecuteActivity(ctx, DestroyInternetGatewayActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

func CreateNatGatewayWorkflow(ctx workflow.Context, input CreateNatGatewayInput) error {
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

	// Create the NatGateway
	if err := workflow.ExecuteActivity(ctx, CreateNatGatewayActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

func DestroyNatGatewayWorkflow(ctx workflow.Context, input DestroyNatGatewayInput) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour,
		HeartbeatTimeout:    time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    10 * time.Second,
			BackoffCoefficient: 1.3,
			MaximumInterval:    2 * time.Minute,
		},
	})

	if err := workflow.ExecuteActivity(ctx, DestroyNatGatewayActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

func CreateRouteTableWorkflow(ctx workflow.Context, input CreateRouteTableInput) error {
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

	// Create the RouteTable
	if err := workflow.ExecuteActivity(ctx, CreateRouteTableActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

func DestroyRouteTableWorkflow(ctx workflow.Context, input DestroyRouteTableInput) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Hour,
		HeartbeatTimeout:    time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    10 * time.Second,
			BackoffCoefficient: 1.3,
			MaximumInterval:    2 * time.Minute,
		},
	})

	if err := workflow.ExecuteActivity(ctx, DestroyRouteTableActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

func CreateRouteTableAssociationWorkflow(ctx workflow.Context, input CreateRouteTableAssociationInput) error {
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

	// Create the RouteTableAssociation
	if err := workflow.ExecuteActivity(ctx, CreateRouteTableAssociationActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

func DestroyRouteTableAssociationWorkflow(ctx workflow.Context, input DestroyRouteTableAssociationInput) error {
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

	// Destroy the RouteTableAssociation
	if err := workflow.ExecuteActivity(ctx, DestroyRouteTableAssociationActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

func CreateRouteWorkflow(ctx workflow.Context, input CreateRouteInput) error {
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

	// Create the Route
	if err := workflow.ExecuteActivity(ctx, CreateRouteActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

func DestroyRouteWorkflow(ctx workflow.Context, input DestroyRouteInput) error {
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

	// Destroy the Route
	if err := workflow.ExecuteActivity(ctx, DestroyRouteActivity, input).Get(ctx, nil); err != nil {
		return err
	}

	return nil
}
