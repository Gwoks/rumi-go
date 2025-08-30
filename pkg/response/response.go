package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo contains error details
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// PaginationMeta contains pagination information
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool            `json:"success"`
	Message    string          `json:"message"`
	Data       interface{}     `json:"data"`
	Pagination *PaginationMeta `json:"pagination"`
	Error      *ErrorInfo      `json:"error,omitempty"`
}

// Success sends a successful response
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a created response
func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// BadRequest sends a bad request error response
func BadRequest(c *gin.Context, message string, details ...string) {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}

	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: "Bad Request",
		Error: &ErrorInfo{
			Code:    "BAD_REQUEST",
			Message: message,
			Details: detail,
		},
	})
}

// Unauthorized sends an unauthorized error response
func Unauthorized(c *gin.Context, message string, details ...string) {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}

	c.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Message: "Unauthorized",
		Error: &ErrorInfo{
			Code:    "UNAUTHORIZED",
			Message: message,
			Details: detail,
		},
	})
}

// Forbidden sends a forbidden error response
func Forbidden(c *gin.Context, message string, details ...string) {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}

	c.JSON(http.StatusForbidden, Response{
		Success: false,
		Message: "Forbidden",
		Error: &ErrorInfo{
			Code:    "FORBIDDEN",
			Message: message,
			Details: detail,
		},
	})
}

// NotFound sends a not found error response
func NotFound(c *gin.Context, message string, details ...string) {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}

	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Message: "Not Found",
		Error: &ErrorInfo{
			Code:    "NOT_FOUND",
			Message: message,
			Details: detail,
		},
	})
}

// Conflict sends a conflict error response
func Conflict(c *gin.Context, message string, details ...string) {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}

	c.JSON(http.StatusConflict, Response{
		Success: false,
		Message: "Conflict",
		Error: &ErrorInfo{
			Code:    "CONFLICT",
			Message: message,
			Details: detail,
		},
	})
}

// InternalServerError sends an internal server error response
func InternalServerError(c *gin.Context, message string, details ...string) {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}

	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Message: "Internal Server Error",
		Error: &ErrorInfo{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: message,
			Details: detail,
		},
	})
}

// ValidationError sends a validation error response
func ValidationError(c *gin.Context, message string, details ...string) {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}

	c.JSON(http.StatusUnprocessableEntity, Response{
		Success: false,
		Message: "Validation Error",
		Error: &ErrorInfo{
			Code:    "VALIDATION_ERROR",
			Message: message,
			Details: detail,
		},
	})
}

// Paginated sends a paginated response
func Paginated(c *gin.Context, message string, data interface{}, pagination *PaginationMeta) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

// CalculatePagination calculates pagination metadata
func CalculatePagination(page, limit int, total int64) *PaginationMeta {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	return &PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}
