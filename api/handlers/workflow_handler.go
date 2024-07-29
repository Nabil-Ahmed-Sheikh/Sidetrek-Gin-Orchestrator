package handlers

import (
	"GinProject/app/worker/workflows"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.uber.org/zap"
)

type (
	WorkflowHandler struct {
		temporalClient client.Client
		taskQueue      string
		logger         *zap.SugaredLogger
	}

	WorkflowInfo struct {
		RunID      string `json:"run_id"`
		WorkflowID string `json:"wf_id"`
	}
)

func NewWorkflowHandler(r *gin.RouterGroup, temporalClient client.Client, taskQueue string, logger *zap.SugaredLogger) *WorkflowHandler {
	handler := &WorkflowHandler{
		temporalClient: temporalClient,
		taskQueue:      taskQueue,
		logger:         logger.Named("workflows_handler"),
	}

	wfs := r.Group("/workflows")
	wfs.POST("/run/:wf", handler.RunInfraAction)

	return handler
}

func (h *WorkflowHandler) RunInfraAction(ctx *gin.Context) {
	wfName := ctx.Param("wf")
	input, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		print(err)
		h.logger.Errorw("Error reading body", zap.Error(err))
		ctx.Status(http.StatusBadRequest)
		return
	}

	in, err := workflows.MapInputToWorkflow(wfName, string(input))
	if err != nil {
		print(err)
		h.logger.Errorw("Error parsing input", zap.String("wfName", wfName), zap.ByteString("input", input), zap.Error(err))
		ctx.Status(http.StatusBadRequest)

		return
	}

	run, err := h.temporalClient.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		TaskQueue:          h.taskQueue,
		WorkflowRunTimeout: time.Second * 6000,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Second * 100,
			MaximumAttempts:    3,
		},
	}, wfName, in)

	if err != nil {
		h.logger.Errorw("Error starting workflow", zap.String("wfName", wfName), zap.ByteString("input", input), zap.Error(err))
		ctx.Status(http.StatusInternalServerError)

		return
	}

	ctx.JSON(http.StatusOK, WorkflowInfo{
		RunID:      run.GetRunID(),
		WorkflowID: run.GetID(),
	})

}
