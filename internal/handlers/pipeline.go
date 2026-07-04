package handlers

import (
	"forgeflow/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PipelineHandler struct {
	pipelineSvc services.PipelineService
	runSvc      services.PipelineRunService
}

func NewPipelineHandler(pipelineSvc services.PipelineService, runSvc services.PipelineRunService) *PipelineHandler {
	return &PipelineHandler{
		pipelineSvc: pipelineSvc,
		runSvc:      runSvc,
	}
}

func (h *PipelineHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("", h.Create)
	// router.GET("", h.List)
	router.POST("/:id/run", h.Run)
}

type CreatePipelineRequest struct {
	ProjectID uuid.UUID `json:"project_id" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	Config    string    `json:"config" binding:"required"`
}

func (h *PipelineHandler) Create(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var req CreatePipelineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	pipeline, err := h.pipelineSvc.CreatePipeline(c.Request.Context(), req.ProjectID, userID, req.Name, req.Config)
	if err != nil {
		if err == services.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, pipeline)
}

func (h *PipelineHandler) Run(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	pipelineID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	run, err := h.runSvc.CreatePipelineRun(c.Request.Context(), pipelineID, userID)
	if err != nil {
		if err == services.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Create execution plan and spawn jobs
	jobs, err := h.runSvc.GenerateExecutionPlan(c.Request.Context(), run)
	if err == nil {
		_ = h.runSvc.CreateJobs(c.Request.Context(), jobs)
		_ = h.runSvc.InitializePipelineStatus(c.Request.Context(), run.ID)
	}

	c.JSON(http.StatusCreated, run)
}
