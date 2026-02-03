package handler

import (
	"strconv"
	"strings"

	"github.com/afghan/rumi-backend/internal/domain/entity"
	"github.com/afghan/rumi-backend/internal/domain/usecase"
	"github.com/afghan/rumi-backend/internal/middleware"
	"github.com/afghan/rumi-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// ChildHandler handles child management related requests
type ChildHandler struct {
	childUsecase usecase.ChildUsecase
}

// NewChildHandler creates a new child handler
func NewChildHandler(childUsecase usecase.ChildUsecase) *ChildHandler {
	return &ChildHandler{
		childUsecase: childUsecase,
	}
}

// CreateChild handles POST /api/v1/children - Create a new child
func (h *ChildHandler) CreateChild(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req entity.CreateChildRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Trim whitespace
	req.Name = strings.TrimSpace(req.Name)
	req.NickName = strings.TrimSpace(req.NickName)
	req.BirthDate = strings.TrimSpace(req.BirthDate)

	child, err := h.childUsecase.CreateChild(c.Request.Context(), userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid birth date") || strings.Contains(err.Error(), "birth date cannot be") {
			response.BadRequest(c, err.Error())
			return
		}
		if strings.Contains(err.Error(), "required") {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalServerError(c, "Failed to create child", err.Error())
		return
	}

	response.Created(c, "Child created successfully", child)
}

// GetChildren handles GET /api/v1/children - Get user's children
func (h *ChildHandler) GetChildren(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	children, total, err := h.childUsecase.GetUserChildren(c.Request.Context(), userID, limit, offset)
	if err != nil {
		response.InternalServerError(c, "Failed to get children", err.Error())
		return
	}

	// Calculate pagination
	pagination := response.CalculatePagination(page, limit, total)

	response.Paginated(c, "Children retrieved successfully", children, pagination)
}

// GetAllChildren handles GET /api/v1/admin/children - Get all children (admin only)
func (h *ChildHandler) GetAllChildren(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	children, total, err := h.childUsecase.GetAllChildren(c.Request.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(c, "Failed to get all children", err.Error())
		return
	}

	// Calculate pagination
	pagination := response.CalculatePagination(page, limit, total)

	response.Paginated(c, "All children retrieved successfully", children, pagination)
}

// GetChild handles GET /api/v1/children/:id - Get child by ID
func (h *ChildHandler) GetChild(c *gin.Context) {
	childID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid child ID")
		return
	}

	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	userRole, exists := middleware.GetUserRoleFromContext(c)
	if !exists {
		response.Unauthorized(c, "User role not found")
		return
	}

	child, err := h.childUsecase.GetChildByID(c.Request.Context(), childID, userID, userRole)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "Child not found")
			return
		}
		response.InternalServerError(c, "Failed to get child", err.Error())
		return
	}

	response.Success(c, "Child retrieved successfully", child)
}

// UpdateChild handles PUT /api/v1/children/:id - Update child
func (h *ChildHandler) UpdateChild(c *gin.Context) {
	childID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid child ID")
		return
	}

	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	userRole, exists := middleware.GetUserRoleFromContext(c)
	if !exists {
		response.Unauthorized(c, "User role not found")
		return
	}

	var req entity.UpdateChildRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Trim whitespace
	req.Name = strings.TrimSpace(req.Name)
	req.NickName = strings.TrimSpace(req.NickName)
	req.BirthDate = strings.TrimSpace(req.BirthDate)

	child, err := h.childUsecase.UpdateChild(c.Request.Context(), childID, userID, userRole, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "Child not found")
			return
		}
		if strings.Contains(err.Error(), "invalid birth date") || strings.Contains(err.Error(), "birth date cannot be") {
			response.BadRequest(c, err.Error())
			return
		}
		if strings.Contains(err.Error(), "required") {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalServerError(c, "Failed to update child", err.Error())
		return
	}

	response.Success(c, "Child updated successfully", child)
}

// DeleteChild handles DELETE /api/v1/children/:id - Delete child
func (h *ChildHandler) DeleteChild(c *gin.Context) {
	childID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid child ID")
		return
	}

	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	userRole, exists := middleware.GetUserRoleFromContext(c)
	if !exists {
		response.Unauthorized(c, "User role not found")
		return
	}

	err = h.childUsecase.DeleteChild(c.Request.Context(), childID, userID, userRole)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "Child not found")
			return
		}
		response.InternalServerError(c, "Failed to delete child", err.Error())
		return
	}

	response.Success(c, "Child deleted successfully", nil)
}

// SetChildActive handles PUT /api/v1/admin/children/:id/active - Set child active status (admin only)
func (h *ChildHandler) SetChildActive(c *gin.Context) {
	childID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid child ID")
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	err = h.childUsecase.SetChildActive(c.Request.Context(), childID, req.IsActive)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "Child not found")
			return
		}
		response.InternalServerError(c, "Failed to update child status", err.Error())
		return
	}

	status := "activated"
	if !req.IsActive {
		status = "deactivated"
	}

	response.Success(c, "Child "+status+" successfully", nil)
}
