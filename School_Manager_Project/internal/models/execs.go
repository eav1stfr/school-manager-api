package models

import "database/sql"

type Exec struct {
	ID                   int            `json:"id" db:"id"`
	FirstName            string         `json:"first_name" db:"first_name" validate:"required"`
	LastName             string         `json:"last_name" db:"last_name" validate:"required"`
	Email                string         `json:"email" db:"email" validate:"required"`
	Username             string         `json:"username" db:"username" validate:"required"`
	Password             string         `json:"password" db:"password" validate:"required"`
	PasswordChangedAt    sql.NullString `json:"password_changed_at" db:"password_changed_at"`
	UserCreatedAt        sql.NullString `json:"user_created_at" db:"user_created_at"`
	PasswordResetToken   sql.NullString `json:"password_reset_token" db:"password_reset_token"`
	PasswordTokenExpires sql.NullString `json:"password_token_expires" db:"password_token_expires"`
	InactiveStatus       bool           `json:"inactive_status" db:"inactive_status"`
	Role                 string         `json:"role" db:"role" validate:"required"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required"`
}
