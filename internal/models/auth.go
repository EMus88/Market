package models

type UpdateRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
