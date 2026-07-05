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
	router.GET("", h.List)
	router.POST("/:id/run", h.Run)
	router.GET("/:id/runs", h.ListRuns)
	router.GET("/:id/runs/:run_id", h.GetRunDetails)
	router.DELETE("/:id", h.Delete)
}

type CreatePipelineRequest struct {
	ProjectID   uuid.UUID `json:"project_id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	YamlConfig  string    `json:"yaml_config" binding:"required"`
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

	pipeline, err := h.pipelineSvc.CreatePipeline(c.Request.Context(), req.ProjectID, userID, req.Name, req.Description, req.YamlConfig)
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

func (h *PipelineHandler) List(c *gin.Context) {
	projectIDStr := c.Query("project_id")
	if projectIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id query param is required"})
		return
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project_id"})
		return
	}

	pipelines, err := h.pipelineSvc.ListPipelinesByProject(c.Request.Context(), projectID, 100, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, pipelines)
}

func (h *PipelineHandler) ListRuns(c *gin.Context) {
	pipelineID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	runs, err := h.runSvc.ListRunsByPipeline(c.Request.Context(), pipelineID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, runs)
}

func (h *PipelineHandler) GetRunDetails(c *gin.Context) {
	runID, err := uuid.Parse(c.Param("run_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid run ID"})
		return
	}

	run, err := h.runSvc.GetRunDetails(c.Request.Context(), runID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, run)
}

func (h *PipelineHandler) Delete(c *gin.Context) {
	pipelineID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	err = h.pipelineSvc.DeletePipeline(c.Request.Context(), pipelineID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pipeline deleted successfully"})
}
