package models

type Teacher struct {
	ID        int    `json:"id" db:"id"`
	FirstName string `json:"first_name" db:"first_name" validate:"required"`
	LastName  string `json:"last_name" db:"last_name" validate:"required"`
	Email     string `json:"email" db:"email" validate:"required,email"`
	Class     string `json:"class" db:"class" validate:"required"`
	Subject   string `json:"subject" db:"subject" validate:"required"`
}
