package network

// func findRolePolicyAttachmentByNameAndArn(ctx context.Context, awsConfig aws.Config, roleName, policyArn string) (types.AttachedPolicy, error) {
// 	client := iam.NewFromConfig(awsConfig)

// 	// TODO: if it has more than 100 policies, we need to paginate
// 	attachOutput, err := client.ListAttachedRolePolicies(ctx, &iam.ListAttachedRolePoliciesInput{
// 		RoleName: aws.String(roleName),
// 	})

// 	if err != nil {
// 		return types.AttachedPolicy{}, err
// 	}

// 	if attachOutput == nil {
// 		return types.AttachedPolicy{}, nil
// 	}

// 	if len(attachOutput.AttachedPolicies) == 0 {
// 		return types.AttachedPolicy{}, nil
// 	}

// 	for _, policy := range attachOutput.AttachedPolicies {
// 		if policy.PolicyArn == nil {
// 			continue
// 		}

// 		if *policy.PolicyArn == policyArn {
// 			return policy, nil
// 		}
// 	}

// 	return types.AttachedPolicy{}, nil
// }
