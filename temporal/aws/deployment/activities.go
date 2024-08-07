package deployment

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

func CreateDeploymentActivity(ctx context.Context, input CreateDeploymentInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	// logger := activity.GetLogger(ctx)

	// Temporal activity aware Terraform workspace wrapper
	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "aws/deployment",
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			// DynamoDB:    cfg.TfState.DynamoDB,
			Key: fmt.Sprintf("deployment-%s.tfstate", input),
		},
	})

	fmt.Println("OOOOOOOOOOOOOOOOOO")
	fmt.Println(input.ClusterName)
	fmt.Println("OOOOOOOOOOOOOOOOOO")

	// Apply Terraform
	_, err := tfa.Apply(ctx, tfworkspace.ApplyInput{
		AwsCredentials: awsConfig.Credentials,
		Env: map[string]string{
			"AWS_REGION": cfg.TfState.Region,
		},
		Vars: map[string]interface{}{
			"cluster_name": input.ClusterName,
			// 	// "ecr_address":     input.EcrAddress,
			// 	// "ecr_user":        user,
			// 	// "ecr_password":    pass,
			// 	// "ecr_repo":        input.EcrRepoName,
			// 	// "source_path":     input.SourcePath,
			// 	// we can add other from `input`
		},
	})

	// Extract output from Terraform
	// Ec2ID, err := applyOutput.String("vpc_id")
	if err != nil {
		return nil
	} else {
		return err
	}

}

func DestroyDeploymentActivity(ctx context.Context, input DestroymentInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "aws/deployment",
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			// DynamoDB:    cfg.TfState.DynamoDB,
			Key: fmt.Sprintf("deployment-%s.tfstate", input),
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
