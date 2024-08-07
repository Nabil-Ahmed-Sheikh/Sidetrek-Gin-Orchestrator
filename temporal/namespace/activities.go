package namespace

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

func CreateNamespaceActivity(ctx context.Context, input CreateNamespaceInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	// Temporal activity aware Terraform workspace wrapper
	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "namespace",
		TerraformFS:   terraform.Namespace,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("namespace-%s.tfstate", input.NamespaceName),
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
			"namespace_name": input.NamespaceName,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func DestroyNamespaceActivity(ctx context.Context, input DestroyNamespaceInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "namespace",
		TerraformFS:   terraform.Namespace,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("namespace-%s.tfstate", input.NamespaceName),
		},
	})

	if err := tfa.Destroy(ctx, tfworkspace.DestroyInput{
		AwsCredentials: awsConfig.Credentials,
		Env: map[string]string{
			"AWS_REGION": cfg.TfState.Region,
		},
		Vars: map[string]interface{}{
			"cluster_name":   input.ClusterName,
			"namespace_name": input.NamespaceName,
		},
	}); err != nil {
		return err
	}
	return nil
}
