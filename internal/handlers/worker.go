package handlers

import (
	"forgeflow/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WorkerHandler struct {
	workerSvc services.WorkerService
}

func NewWorkerHandler(workerSvc services.WorkerService) *WorkerHandler {
	return &WorkerHandler{workerSvc: workerSvc}
}

func (h *WorkerHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/register", h.Register)
	router.POST("/:id/heartbeat", h.Heartbeat)
}

type RegisterWorkerRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *WorkerHandler) Register(c *gin.Context) {
	var req RegisterWorkerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	worker, err := h.workerSvc.RegisterWorker(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, worker)
}

func (h *WorkerHandler) Heartbeat(c *gin.Context) {
	workerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid worker ID"})
		return
	}

	if err := h.workerSvc.ReceiveHeartbeat(c.Request.Context(), workerID); err != nil {
		if err == services.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Worker not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Heartbeat received"})
}
