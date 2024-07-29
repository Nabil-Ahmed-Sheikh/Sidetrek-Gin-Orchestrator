package workflows

import (
	"os/exec"

	"go.temporal.io/sdk/workflow"
)

func TerraformWorkflow(ctx workflow.Context) error {
	// Initialize Terraform
	initCmd := exec.Command("terraform", "init")
	initCmd.Dir = "./terraform"
	if err := initCmd.Run(); err != nil {
		return err
	}

	// Apply Terraform
	applyCmd := exec.Command("terraform", "apply", "-auto-approve")
	applyCmd.Dir = "./terraform"
	if err := applyCmd.Run(); err != nil {
		return err
	}

	return nil
}
