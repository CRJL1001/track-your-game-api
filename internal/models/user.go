package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Email          string    `json:"email" db:"email" validate:"required"`
	HashedPassword string    `json:"-" db:"password"` //- pour exclure du json
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`    //format email requis
	Password string `json:"password" validate:"required,min=8"` //minimum 8 caractères
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"` //format email requis
	Password string `json:"password" validate:"required"`    //minimum 8 caractères
}
