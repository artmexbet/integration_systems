package v1

import "time"

// ErrorResponse represents an API error response
// @Description Error response with message
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// StatsResponse represents statistics about the dataset
// @Description Dataset statistics
type StatsResponse struct {
	LaureatesCount  int64 `json:"laureates_count"`
	PrizesCount     int64 `json:"prizes_count"`
	CategoriesCount int64 `json:"categories_count"`
}

// LastUpdateResponse represents the last update timestamp
// @Description Last dataset update information
type LastUpdateResponse struct {
	LastUpdate time.Time `json:"last_update"`
}

// LaureateResponse represents a laureate in API responses
// @Description Nobel laureate information
type LaureateResponse struct {
	ID         int32   `json:"id"`
	Firstname  string  `json:"firstname"`
	Surname    string  `json:"surname,omitempty"`
	Motivation string  `json:"motivation"`
	Share      int32   `json:"share"`
	UpdatedAt  *string `json:"updated_at,omitempty"`
}

// LaureateListResponse represents a list of laureates
// @Description List of laureates with pagination info
type LaureateListResponse struct {
	Data       []LaureateResponse `json:"data"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PerPage    int                `json:"per_page"`
	TotalPages int                `json:"total_pages"`
}

// CreateLaureateRequest represents the request to create a laureate
// @Description Create laureate request body
type CreateLaureateRequest struct {
	ID         int32  `json:"id" validate:"required"`
	Firstname  string `json:"firstname" validate:"required"`
	Surname    string `json:"surname,omitempty"`
	Motivation string `json:"motivation" validate:"required"`
	Share      int32  `json:"share" validate:"required,min=1,max=4"`
}

// UpdateLaureateRequest represents the request to update a laureate
// @Description Update laureate request body
type UpdateLaureateRequest struct {
	Firstname  string `json:"firstname" validate:"required"`
	Surname    string `json:"surname,omitempty"`
	Motivation string `json:"motivation" validate:"required"`
	Share      int32  `json:"share" validate:"required,min=1,max=4"`
}

// PrizeResponse represents a prize in API responses
// @Description Nobel prize information
type PrizeResponse struct {
	ID        int32              `json:"id"`
	Year      int32              `json:"year"`
	Category  string             `json:"category"`
	Laureates []LaureateResponse `json:"laureates,omitempty"`
	UpdatedAt *string            `json:"updated_at,omitempty"`
}

// PrizeListResponse represents a list of prizes
// @Description List of prizes with pagination info
type PrizeListResponse struct {
	Data       []PrizeResponse `json:"data"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PerPage    int             `json:"per_page"`
	TotalPages int             `json:"total_pages"`
}

// CreatePrizeRequest represents the request to create a prize
// @Description Create prize request body
type CreatePrizeRequest struct {
	Year        int32   `json:"year" validate:"required,min=1901"`
	Category    string  `json:"category" validate:"required"`
	LaureateIDs []int32 `json:"laureate_ids,omitempty"`
}

// UpdatePrizeRequest represents the request to update a prize
// @Description Update prize request body
type UpdatePrizeRequest struct {
	Year     int32  `json:"year" validate:"required,min=1901"`
	Category string `json:"category" validate:"required"`
}

// CategoriesResponse represents a list of categories
// @Description List of prize categories
type CategoriesResponse struct {
	Categories []string `json:"categories"`
}

// SuccessResponse represents a generic success response
// @Description Generic success response
type SuccessResponse struct {
	Message string `json:"message"`
}
