package handlers

import (
	"forgeflow/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type JobHandler struct {
	jobSvc services.JobService
}

func NewJobHandler(jobSvc services.JobService) *JobHandler {
	return &JobHandler{jobSvc: jobSvc}
}

func (h *JobHandler) RegisterRoutes(router *gin.RouterGroup) {
	// router.GET("", h.List)
	router.POST("/:id/retry", h.Retry)
}

func (h *JobHandler) Retry(c *gin.Context) {
	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	err = h.jobSvc.RetryFailedJob(c.Request.Context(), jobID)
	if err != nil {
		if err == services.ErrInvalidStateTransition {
			c.JSON(http.StatusConflict, gin.H{"error": "Job is not in a failed state"})
			return
		}
		if err == services.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Job queued for retry"})
}
