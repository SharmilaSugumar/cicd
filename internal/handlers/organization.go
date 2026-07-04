package handlers

import (
	"forgeflow/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrganizationHandler struct {
	orgSvc services.OrganizationService
}

func NewOrganizationHandler(orgSvc services.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{orgSvc: orgSvc}
}

func (h *OrganizationHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("", h.Create)
	router.GET("", h.List)
	// Additional routes for members can be added here
}

type CreateOrgRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *OrganizationHandler) Create(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var req CreateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	org, err := h.orgSvc.CreateOrganization(c.Request.Context(), req.Name, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":   org.ID,
		"name": org.Name,
	})
}

func (h *OrganizationHandler) List(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	orgs, err := h.orgSvc.ListOrganizationsForUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Minimal response
	c.JSON(http.StatusOK, orgs)
}
