package domain

type Device struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}
