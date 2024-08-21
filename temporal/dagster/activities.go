package dagster

import (
	"context"
	"fmt"

	"GinProject/app/config/awsconfig"
	"GinProject/app/config/env"
	"GinProject/app/orchestrator/tfactivity"
	"GinProject/app/orchestrator/tfexec"
	"GinProject/app/orchestrator/tfworkspace"
	"GinProject/app/terraform"
)

func CreateDagsterClusterActivity(ctx context.Context, input CreateDagsterClusterInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	// Temporal activity aware Terraform workspace wrapper
	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "dagster",
		TerraformFS:   terraform.Dagster,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("dagster-cluster-%s.tfstate", input.ClusterName),
		},
	})

	// Apply Terraform
	_, err := tfa.Apply(ctx, tfworkspace.ApplyInput{
		AwsCredentials: awsConfig.Credentials,
		Env: map[string]string{
			"AWS_REGION": cfg.TfState.Region,
		},
		Vars: map[string]interface{}{
			"cluster_name":   input.ClusterName,
			"repository":     input.repository,
			"additional_set": input.AdditionalSet,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func DestroyDagsterClusterActivity(ctx context.Context, input DestroyDagsterClusterInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "dagster",
		TerraformFS:   terraform.Dagster,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("dagster-cluster-%s.tfstate", input.ClusterName),
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
