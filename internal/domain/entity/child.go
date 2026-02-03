package entity

import (
	"time"
)

// Child represents a child in the system
type Child struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	NickName  string    `json:"nick_name" db:"nick_name"`
	BirthDate time.Time `json:"birth_date" db:"birth_date"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Join fields
	User *PublicUser `json:"user,omitempty" db:"-"`
}

// CreateChildRequest represents request to create a new child
type CreateChildRequest struct {
	Name      string `json:"name" binding:"required,min=2,max=100"`
	NickName  string `json:"nick_name" binding:"max=50"`
	BirthDate string `json:"birth_date" binding:"required"` // Format: YYYY-MM-DD
}

// UpdateChildRequest represents request to update a child
type UpdateChildRequest struct {
	Name      string `json:"name" binding:"required,min=2,max=100"`
	NickName  string `json:"nick_name" binding:"max=50"`
	BirthDate string `json:"birth_date" binding:"required"` // Format: YYYY-MM-DD
	IsActive  *bool  `json:"is_active,omitempty"`
}

// ChildListResponse represents response for child list
type ChildListResponse struct {
	Children   []*Child        `json:"children"`
	Pagination *PaginationMeta `json:"pagination"`
}
