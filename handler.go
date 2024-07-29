package app

import (
	"context"
	"net/http"

	"GinProject/app/workflows"

	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
)

// @@@SNIPSTART applyTerraform
func ApplyTerraform(gc *gin.Context) {
	c, err := client.Dial(client.Options{})
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create Temporal client"})
		return
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "terraform_workflow",
		TaskQueue: "terraform-namespace",
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, workflows.TerraformWorkflow)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to execute workflow"})
		return
	}

	gc.JSON(http.StatusOK, gin.H{"workflowID": we.GetID(), "runID": we.GetRunID()})
}

// @@@SNIPEND
