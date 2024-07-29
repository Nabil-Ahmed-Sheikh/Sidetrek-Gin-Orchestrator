package ec2

import (
	"GinProject/app/config/awsconfig"
	"GinProject/app/config/env"
	"GinProject/app/orchestrator/tfactivity"
	"GinProject/app/orchestrator/tfexec"
	"GinProject/app/orchestrator/tfworkspace"
	"GinProject/app/terraform"
	"context"
	"fmt"
)

func CreateEc2Activity(ctx context.Context, input CreateEc2Input) (CreateEc2Output, error) {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	// logger := activity.GetLogger(ctx)

	// Temporal activity aware Terraform workspace wrapper
	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "aws/ec2",
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Key:         fmt.Sprintf("ecr-%s.tfstate", input),
		},
	})

	// Apply Terraform
	_, err := tfa.Apply(ctx, tfworkspace.ApplyInput{
		AwsCredentials: awsConfig.Credentials,
		Env: map[string]string{
			"AWS_REGION": cfg.TfState.Region,
		},
		Vars: map[string]interface{}{
			// "create_ecr_repo": input.Name,
			// "ecr_address":     input.EcrAddress,
			// "ecr_user":        user,
			// "ecr_password":    pass,
			// "ecr_repo":        input.EcrRepoName,
			// "source_path":     input.SourcePath,
			// we can add other from `input`
		},
	})

	// Extract output from Terraform
	// Ec2ID, err := applyOutput.String("vpc_id")
	if err != nil {
		return CreateEc2Output{}, err
	}

	return CreateEc2Output{
		// Ec2ID: Ec2ID,
	}, nil

}

func DestroyEc2Activity(ctx context.Context, input DestroyEc2Output) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "aws/vpc",
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			// Bucket:      cfg.TfState.Bucket,
			// DynamoDB:    cfg.TfState.DynamoDB,
			Key: fmt.Sprintf("vpc-%s.tfstate", input),
		},
	})

	if err := tfa.Destroy(ctx, tfworkspace.DestroyInput{
		AwsCredentials: awsConfig.Credentials,
		Env: map[string]string{
			"AWS_REGION": cfg.TfState.Region,
		},
	}); err != nil {
		return err
	}
	return nil
}
