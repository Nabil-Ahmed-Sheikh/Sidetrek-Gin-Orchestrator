package emissary

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

func CreateEmissaryMappingActivity(ctx context.Context, input CreateEmissaryMappingInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	// Temporal activity aware Terraform workspace wrapper
	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "emissary",
		TerraformFS:   terraform.Emissary,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("emissary-mapping-%s.tfstate", input.HostName),
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
			"svc_name":       input.ServiceName,
			"host_name":      input.HostName,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func DestroyEmissaryMappingActivity(ctx context.Context, input DestroyEmissaryMappingInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "emissary",
		TerraformFS:   terraform.Emissary,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("emissary-mapping-%s.tfstate", input.HostName),
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
			"svc_name":       input.ServiceName,
			"host_name":      input.HostName,
		},
	}); err != nil {
		return err
	}

	return nil
}
