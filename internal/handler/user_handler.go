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

// UserHandler handles user management related requests
type UserHandler struct {
	userManagementUsecase usecase.UserManagementUsecase
	userActivityUsecase   usecase.UserActivityUsecase
}

// NewUserHandler creates a new user handler
func NewUserHandler(
	userManagementUsecase usecase.UserManagementUsecase,
	userActivityUsecase usecase.UserActivityUsecase,
) *UserHandler {
	return &UserHandler{
		userManagementUsecase: userManagementUsecase,
		userActivityUsecase:   userActivityUsecase,
	}
}

// GetAllUsers handles GET /api/v1/users - Get all users (admin only)
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	users, total, err := h.userManagementUsecase.GetAllUsers(c.Request.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(c, "Failed to get users", err.Error())
		return
	}

	// Calculate pagination
	pagination := response.CalculatePagination(page, limit, total)

	response.Paginated(c, "Users retrieved successfully", users, pagination)
}

// GetUser handles GET /api/v1/users/:id - Get user by ID (admin only)
func (h *UserHandler) GetUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	user, err := h.userManagementUsecase.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "Failed to get user", err.Error())
		return
	}

	response.Success(c, "User retrieved successfully", user)
}

// UpdateUser handles PUT /api/v1/users/:id - Update user (admin only)
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req entity.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Normalize email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Name = strings.TrimSpace(req.Name)
	req.Phone = strings.TrimSpace(req.Phone)

	user, err := h.userManagementUsecase.UpdateUser(c.Request.Context(), userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "User not found")
			return
		}
		if strings.Contains(err.Error(), "email already exists") {
			response.Conflict(c, "Email address is already registered")
			return
		}
		if strings.Contains(err.Error(), "password must be") || strings.Contains(err.Error(), "invalid") {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalServerError(c, "Failed to update user", err.Error())
		return
	}

	response.Success(c, "User updated successfully", user)
}

// DeleteUser handles DELETE /api/v1/users/:id - Delete user (admin only)
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	err = h.userManagementUsecase.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "Failed to delete user", err.Error())
		return
	}

	response.Success(c, "User deleted successfully", nil)
}

// SetUserActive handles PUT /api/v1/users/:id/active - Set user active status (admin only)
func (h *UserHandler) SetUserActive(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	err = h.userManagementUsecase.SetUserActive(c.Request.Context(), userID, req.IsActive)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "Failed to update user status", err.Error())
		return
	}

	status := "activated"
	if !req.IsActive {
		status = "deactivated"
	}

	response.Success(c, "User "+status+" successfully", nil)
}

// UpdateProfile handles PUT /api/v1/profile - Update user's own profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req entity.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Trim whitespace
	req.Name = strings.TrimSpace(req.Name)
	req.Phone = strings.TrimSpace(req.Phone)

	user, err := h.userManagementUsecase.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "Failed to update profile", err.Error())
		return
	}

	response.Success(c, "Profile updated successfully", user)
}

// ChangePassword handles PUT /api/v1/profile/password - Change user's password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req entity.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	err := h.userManagementUsecase.ChangePassword(c.Request.Context(), userID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "User not found")
			return
		}
		if strings.Contains(err.Error(), "incorrect") {
			response.BadRequest(c, "Current password is incorrect")
			return
		}
		if strings.Contains(err.Error(), "password must be") {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalServerError(c, "Failed to change password", err.Error())
		return
	}

	response.Success(c, "Password changed successfully", nil)
}

// GetUserActivities handles GET /api/v1/users/:id/activities - Get user activities (admin) or own activities
func (h *UserHandler) GetUserActivities(c *gin.Context) {
	targetUserID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	currentUserID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	currentUserRole, exists := middleware.GetUserRoleFromContext(c)
	if !exists {
		response.Unauthorized(c, "User role not found")
		return
	}

	// Users can only see their own activities unless they're admin
	if currentUserRole != entity.RoleAdmin && currentUserID != targetUserID {
		response.Forbidden(c, "Cannot view other user's activities")
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	activities, total, err := h.userActivityUsecase.GetUserActivities(c.Request.Context(), targetUserID, limit, offset)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "Failed to get user activities", err.Error())
		return
	}

	// Calculate pagination
	pagination := response.CalculatePagination(page, limit, total)

	response.Paginated(c, "User activities retrieved successfully", activities, pagination)
}

// GetAllActivities handles GET /api/v1/activities - Get all activities (admin only)
func (h *UserHandler) GetAllActivities(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	activities, total, err := h.userActivityUsecase.GetAllActivities(c.Request.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(c, "Failed to get activities", err.Error())
		return
	}

	// Calculate pagination
	pagination := response.CalculatePagination(page, limit, total)

	response.Paginated(c, "Activities retrieved successfully", activities, pagination)
}
