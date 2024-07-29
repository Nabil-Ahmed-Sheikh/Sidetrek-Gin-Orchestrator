package network

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

func CreateVpcActivity(ctx context.Context, input CreateVpcInput) (CreateVpcOutput, error) {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	// Temporal activity aware Terraform workspace wrapper
	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "aws/vpc",
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("vpc-%s.tfstate", input.Name),
		},
	})

	// Apply Terraform
	applyOutput, err := tfa.Apply(ctx, tfworkspace.ApplyInput{
		AwsCredentials: awsConfig.Credentials,
		Env: map[string]string{
			"AWS_REGION": cfg.TfState.Region,
		},
		Vars: map[string]interface{}{
			"cidr_block": input.CIDRBlock,
			"name":       input.Name,
		},
	})
	if err != nil {
		return CreateVpcOutput{}, err
	}

	// Extract output from Terraform
	vpcID, err := applyOutput.String("vpc_id")
	if err != nil {
		return CreateVpcOutput{}, err
	}

	return CreateVpcOutput{
		VpcID: vpcID,
	}, nil
}

func DestroyVpcActivity(ctx context.Context, input DestroyVpcInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "aws/vpc",
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("vpc-%s.tfstate", input.Name),
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

func CreateSubnetsActivity(ctx context.Context, input CreateSubnetsInput) (CreateSubnetsOutput, error) {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	// Temporal activity aware Terraform workspace wrapper
	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "aws/subnet",
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("subnets-%s-%s.tfstate", input.VpcID, input.Subnets.String()),
		},
	})

	var subnets []map[string]string
	for _, s := range input.Subnets {
		subnets = append(subnets, map[string]string{
			"cidr_block":              s.CIDRBlock,
			"name":                    fmt.Sprintf("%s-%s", input.VpcID, s.AvailabilityZone),
			"availability_zone":       cfg.TfState.Region + s.AvailabilityZone,
			"map_public_ip_on_launch": fmt.Sprintf("%t", s.Public),
		})
	}

	// Apply Terraform to create subnets
	if _, err := tfa.Apply(ctx, tfworkspace.ApplyInput{
		AwsCredentials: awsConfig.Credentials,
		Env: map[string]string{
			"AWS_REGION": cfg.TfState.Region,
		},
		Vars: map[string]interface{}{
			"vpc_id":  input.VpcID,
			"subnets": subnets,
		},
	}); err != nil {
		return CreateSubnetsOutput{}, err
	}

	return CreateSubnetsOutput{}, nil
}

func DestroySubnetsActivity(ctx context.Context, input DestroySubnetsInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "aws/subnet",
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("subnets-%s.tfstate", input.VpcID),
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

func CreateSecurityGroupActivity(ctx context.Context, input CreateSecurityGroupInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	// Temporal activity aware Terraform workspace wrapper
	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "aws/security_group",
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("security-group-%s.tfstate", input.Name),
		},
	})

	// Apply Terraform
	_, err := tfa.Apply(ctx, tfworkspace.ApplyInput{
		AwsCredentials: awsConfig.Credentials,
		Env: map[string]string{
			"AWS_REGION": cfg.TfState.Region,
		},
		Vars: map[string]interface{}{
			"name":        input.Name,
			"description": input.Description,
			"vpc_id":      input.VpcID,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func DestroySecurityGroupActivity(ctx context.Context, input DestroySecurityGroupInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	// Temporal activity aware Terraform workspace wrapper
	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "aws/route_table",
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("security-group-%s.tfstate", input.Name),
		},
	})

	// Destroy Terraform
	err := tfa.Destroy(ctx, tfworkspace.DestroyInput{
		AwsCredentials: awsConfig.Credentials,
		Env: map[string]string{
			"AWS_REGION": cfg.TfState.Region,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func CreateElasticIpActivity(ctx context.Context, input CreateElasticIpInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	// Temporal activity aware Terraform workspace wrapper
	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "aws/elastic_ip",
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("elastic-ip-%s.tfstate", input.InternetGatewayID),
		},
	})

	// Apply Terraform
	_, err := tfa.Apply(ctx, tfworkspace.ApplyInput{
		AwsCredentials: awsConfig.Credentials,
		Env: map[string]string{
			"AWS_REGION": cfg.TfState.Region,
		},
		Vars: map[string]interface{}{
			"internet_gateway_id": input.InternetGatewayID,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func DestroyElasticIpActivity(ctx context.Context, input DestroyElasticIpInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "aws/elastic_ip",
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("elastic-ip-%s.tfstate", input.InternetGatewayID),
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

func CreateInternetGatewayActivity(ctx context.Context, input CreateInternetGatewayInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	// Temporal activity aware Terraform workspace wrapper
	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "aws/internet_gateway",
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("internet-gateway-%s.tfstate", input.VpcID),
		},
	})

	// Apply Terraform
	_, err := tfa.Apply(ctx, tfworkspace.ApplyInput{
		AwsCredentials: awsConfig.Credentials,
		Env: map[string]string{
			"AWS_REGION": cfg.TfState.Region,
		},
		Vars: map[string]interface{}{
			"vpc_id": input.VpcID,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func DestroyInternetGatewayActivity(ctx context.Context, input DestroyInternetGatewayInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "aws/internet_gateway",
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("internet-gateway-%s.tfstate", input.VpcID),
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

func CreateNatGatewayActivity(ctx context.Context, input CreateNatGatewayInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	// Temporal activity aware Terraform workspace wrapper
	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "aws/nat_gateway",
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("nat-gateway-%s-%s.tfstate", input.AllocationID, input.SubnetID),
		},
	})

	// Apply Terraform
	_, err := tfa.Apply(ctx, tfworkspace.ApplyInput{
		AwsCredentials: awsConfig.Credentials,
		Env: map[string]string{
			"AWS_REGION": cfg.TfState.Region,
		},
		Vars: map[string]interface{}{
			"nat_eip_id": input.AllocationID,
			"subnet_id":  input.SubnetID,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func DestroyNatGatewayActivity(ctx context.Context, input DestroyNatGatewayInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "aws/nat_gateway",
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("nat-gateway-%s-%s.tfstate", input.AllocationID, input.SubnetID),
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

func CreateRouteTableActivity(ctx context.Context, input CreateRouteTableInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	// Temporal activity aware Terraform workspace wrapper
	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "aws/route_table",
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("route-table-%s-%s.tfstate", input.VpcID, input.Name),
		},
	})

	// Apply Terraform
	_, err := tfa.Apply(ctx, tfworkspace.ApplyInput{
		AwsCredentials: awsConfig.Credentials,
		Env: map[string]string{
			"AWS_REGION": cfg.TfState.Region,
		},
		Vars: map[string]interface{}{
			"vpc_id": input.VpcID,
			"name":   input.Name,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func DestroyRouteTableActivity(ctx context.Context, input DestroyRouteTableInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "aws/route_table",
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("route-table-%s-%s.tfstate", input.VpcID, input.Name),
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

func CreateRouteTableAssociationActivity(ctx context.Context, input CreateRouteTableAssociationInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	// Temporal activity aware Terraform workspace wrapper
	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "aws/route_table_association",
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("route-table-association-%s-%s.tfstate", input.SubnetID, input.RouteTableID),
		},
	})

	// Apply Terraform
	_, err := tfa.Apply(ctx, tfworkspace.ApplyInput{
		AwsCredentials: awsConfig.Credentials,
		Env: map[string]string{
			"AWS_REGION": cfg.TfState.Region,
		},
		Vars: map[string]interface{}{
			"subnet_id":      input.SubnetID,
			"route_table_id": input.RouteTableID,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func DestroyRouteTableAssociationActivity(ctx context.Context, input DestroyRouteTableAssociationInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	// Temporal activity aware Terraform workspace wrapper
	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: "aws/route_table_association",
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         fmt.Sprintf("route-table-association-%s-%s.tfstate", input.SubnetID, input.RouteTableID),
		},
	})

	// Destroy Terraform
	err := tfa.Destroy(ctx, tfworkspace.DestroyInput{
		AwsCredentials: awsConfig.Credentials,
		Env: map[string]string{
			"AWS_REGION": cfg.TfState.Region,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func CreateRouteActivity(ctx context.Context, input CreateRouteInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	// Temporal activity aware Terraform workspace wrapper
	var key string
	var terraformPath string
	if input.NatGatewayID != "" {
		terraformPath = "aws/route_nat"
		key = fmt.Sprintf("route-%s-%s.tfstate", input.RouteTableID, input.NatGatewayID)
	} else {
		terraformPath = "aws/route_igw"
		key = fmt.Sprintf("route-%s-%s.tfstate", input.RouteTableID, input.InternetGatewayID)
	}

	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: terraformPath,
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         key,
		},
	})

	// Apply Terraform
	_, err := tfa.Apply(ctx, tfworkspace.ApplyInput{
		AwsCredentials: awsConfig.Credentials,
		Env: map[string]string{
			"AWS_REGION": cfg.TfState.Region,
		},
		Vars: map[string]interface{}{
			"route_table_id": input.RouteTableID,
			"nat_gateway_id": input.NatGatewayID,
			"gateway_id":     input.InternetGatewayID,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func DestroyRouteActivity(ctx context.Context, input DestroyRouteInput) error {
	cfg := env.MustGetConfig()
	awsConfig := awsconfig.LoadConfig(cfg)

	// Temporal activity aware Terraform workspace wrapper
	var key string
	var terraformPath string
	if input.NatGatewayID != "" {
		terraformPath = "aws/route_nat"
		key = fmt.Sprintf("route-%s-%s.tfstate", input.RouteTableID, input.NatGatewayID)
	} else {
		terraformPath = "aws/route_igw"
		key = fmt.Sprintf("route-%s-%s.tfstate", input.RouteTableID, input.InternetGatewayID)
	}

	tfa := tfactivity.New(tfworkspace.Config{
		TerraformPath: terraformPath,
		TerraformFS:   terraform.AWS,
		Backend: tfexec.BackendConfig{
			Credentials: awsConfig.Credentials,
			Region:      cfg.TfState.Region,
			Bucket:      cfg.TfState.Bucket,
			DynamoDB:    cfg.TfState.DynamoDB,
			Key:         key,
		},
	})

	// Destroy Terraform
	err := tfa.Destroy(ctx, tfworkspace.DestroyInput{
		AwsCredentials: awsConfig.Credentials,
		Env: map[string]string{
			"AWS_REGION": cfg.TfState.Region,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
