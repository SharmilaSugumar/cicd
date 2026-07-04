package handlers

import (
	"forgeflow/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	projectSvc services.ProjectService
}

func NewProjectHandler(projectSvc services.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectSvc: projectSvc}
}

func (h *ProjectHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("", h.Create)
	router.GET("", h.List)
	router.DELETE("/:id", h.Delete)
}

type CreateProjectRequest struct {
	Name           string    `json:"name" binding:"required"`
	Description    string    `json:"description"`
	OrganizationID uuid.UUID `json:"organization_id" binding:"required"`
}

func (h *ProjectHandler) Create(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	project, err := h.projectSvc.CreateProject(c.Request.Context(), req.Name, req.Description, req.OrganizationID, userID)
	if err != nil {
		if err == services.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, project)
}

func (h *ProjectHandler) List(c *gin.Context) {
	orgIDStr := c.Query("organization_id")
	if orgIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organization_id query param is required"})
		return
	}

	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization_id"})
		return
	}

	projects, err := h.projectSvc.ListProjectsByOrganization(c.Request.Context(), orgID, 100, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	err = h.projectSvc.SoftDeleteProject(c.Request.Context(), projectID, userID)
	if err != nil {
		if err == services.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}
		if err == services.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
