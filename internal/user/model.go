package user

import (
	"time"
)

type User struct {
	ID        int
	UserName  string
	Email     string
	Age       int
	CreatedAt time.Time
}

type CreateUserRequest struct {
	UserName string `json:"user_name" validate:"required,min=3,max=50,username_chars"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"omitempty,gte=18"`
}

type UpdateUserRequest struct {
	Email string `json:"email" validate:"omitempty,email"`
	Age   int    `json:"age" validate:"omitempty,gte=18"`
}

type UserResponse struct {
	Message string `json:"msg"`
	Data    []User `json:"data"`
}
