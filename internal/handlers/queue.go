package handlers

import (
	"forgeflow/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type QueueHandler struct {
	queueSvc services.QueueService
}

func NewQueueHandler(queueSvc services.QueueService) *QueueHandler {
	return &QueueHandler{queueSvc: queueSvc}
}

func (h *QueueHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/:id/pause", h.Pause)
	router.POST("/:id/resume", h.Resume)
	router.GET("/:id/stats", h.Stats)
}

func (h *QueueHandler) Pause(c *gin.Context) {
	queueID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid queue ID"})
		return
	}

	if err := h.queueSvc.PauseQueue(c.Request.Context(), queueID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Queue paused"})
}

func (h *QueueHandler) Resume(c *gin.Context) {
	queueID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid queue ID"})
		return
	}

	if err := h.queueSvc.ResumeQueue(c.Request.Context(), queueID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Queue resumed"})
}

func (h *QueueHandler) Stats(c *gin.Context) {
	queueID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid queue ID"})
		return
	}

	stats, err := h.queueSvc.GetQueueStatistics(c.Request.Context(), queueID)
	if err != nil {
		if err == services.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Queue not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
