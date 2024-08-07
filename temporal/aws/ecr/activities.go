package ecr

import (
	"GinProject/app/config/awsconfig"
	"GinProject/app/config/env"
	"GinProject/app/orchestrator/tfactivity"
	"GinProject/app/orchestrator/tfexec"
	"GinProject/app/orchestrator/tfworkspace"
	"GinProject/app/terraform"
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"go.temporal.io/sdk/activity"
)

func ApplyDockerBuildAndPushEcrActivity(ctx context.Context, input ApplyDockerBuildAndPushEcrInput) (ApplyDockerBuildAndPushEcrOutput, error) {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	logger := activity.GetLogger(ctx)

	// get ecr token
	ecrClient := ecr.NewFromConfig(awsConfig)
	token, err := ecrClient.GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		logger.Info("error", err)
		return ApplyDockerBuildAndPushEcrOutput{}, err
	}

	if len(token.AuthorizationData) == 0 || token.AuthorizationData[0].AuthorizationToken == nil {
		return ApplyDockerBuildAndPushEcrOutput{}, fmt.Errorf("no authorization data found")
	}

	decoded, err := base64.StdEncoding.DecodeString(*token.AuthorizationData[0].AuthorizationToken)
	if err != nil {
		return ApplyDockerBuildAndPushEcrOutput{}, err
	}

	authToken := strings.Split(string(decoded), ":")
	fmt.Println("authToken", authToken)

	if len(authToken) != 2 {
		// user:pass
		return ApplyDockerBuildAndPushEcrOutput{}, fmt.Errorf("invalid authorization token")
	}

	user := authToken[0]
	pass := authToken[1]

	// Temporal activity aware Terraform workspace wrapper
	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "aws/ecr",
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("ecr-%s.tfstate", input.EcrRepoName),
		},
	})

	// Apply Terraform
	applyOutput, err := tfa.Apply(ctx, tfworkspace.ApplyInput{
		AwsCredentials: awsConfig.Credentials,
		Env: map[string]string{
			"AWS_REGION": cfg.TfState.Region,
		},
		Vars: map[string]interface{}{
			"create_ecr_repo": input.CreateEcrRepo,
			"ecr_address":     input.EcrAddress,
			"ecr_user":        user,
			"ecr_password":    pass,
			"ecr_repo":        input.EcrRepoName,
			"source_path":     input.SourcePath,
			// we can add other from `input`
		},
	})
	if err != nil {
		return ApplyDockerBuildAndPushEcrOutput{}, err
	}

	// Extract output from Terraform
	imageUri, err := applyOutput.String("image_uri")
	if err != nil {
		return ApplyDockerBuildAndPushEcrOutput{}, err
	}

	imageId, err := applyOutput.String("image_id")
	if err != nil {
		return ApplyDockerBuildAndPushEcrOutput{}, err
	}

	return ApplyDockerBuildAndPushEcrOutput{
		ImageUri: imageUri,
		ImageId:  imageId,
	}, nil
}

func DestroyDockerBuildAndPushEcrActivity(ctx context.Context, input DestroyDockerBuildAndPushEcrInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "aws/ecr",
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("ecr-%s.tfstate", input.EcrRepoName),
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
