package entity

import (
	"encoding/json"
	"time"
)

// ActivityType represents the type of user activity
type ActivityType string

const (
	ActivityAccountCreated  ActivityType = "account_created"
	ActivityLogin           ActivityType = "login"
	ActivityLogout          ActivityType = "logout"
	ActivityProfileUpdate   ActivityType = "profile_update"
	ActivityPasswordChange  ActivityType = "password_change"
	ActivityUserCreated     ActivityType = "user_created"
	ActivityUserUpdated     ActivityType = "user_updated"
	ActivityUserDeleted     ActivityType = "user_deleted"
	ActivityUserActivated   ActivityType = "user_activated"
	ActivityUserDeactivated ActivityType = "user_deactivated"
)

// UserActivity represents a user activity record
type UserActivity struct {
	ID           int64                  `json:"id" db:"id"`
	UserID       int64                  `json:"user_id" db:"user_id"`
	ActivityType ActivityType           `json:"activity_type" db:"activity_type"`
	Description  string                 `json:"description" db:"description"`
	IPAddress    string                 `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent    string                 `json:"user_agent,omitempty" db:"user_agent"`
	Metadata     map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`

	// Join fields
	User *PublicUser `json:"user,omitempty" db:"-"`
}

// CreateActivityRequest represents request to log a user activity
type CreateActivityRequest struct {
	UserID       int64                  `json:"user_id" binding:"required"`
	ActivityType ActivityType           `json:"activity_type" binding:"required"`
	Description  string                 `json:"description" binding:"required"`
	IPAddress    string                 `json:"ip_address,omitempty"`
	UserAgent    string                 `json:"user_agent,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=100"`
	Phone string `json:"phone" binding:"required,min=10,max=15"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

// ActivityListResponse represents response for activity list
type ActivityListResponse struct {
	Activities []*UserActivity `json:"activities"`
	Pagination *PaginationMeta `json:"pagination"`
}

// PaginationMeta represents pagination information
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// MarshalMetadata converts metadata map to JSON string for database storage
func (ua *UserActivity) MarshalMetadata() ([]byte, error) {
	if ua.Metadata == nil {
		return nil, nil
	}
	return json.Marshal(ua.Metadata)
}

// UnmarshalMetadata converts JSON string from database to metadata map
func (ua *UserActivity) UnmarshalMetadata(data []byte) error {
	if len(data) == 0 {
		ua.Metadata = nil
		return nil
	}
	return json.Unmarshal(data, &ua.Metadata)
}
