package user

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int64
	UserName     string
	Email        string
	Age          sql.NullInt64
	CreatedAt    time.Time
}

type CreateUserRequest struct {
	UserName string `json:"user_name" validate:"required,min=3,max=50,username_chars"`
	Email string `json:"email" validate:"required,email"`
	Age  sql.NullInt64 `json:"age" validate:"omitempty,gte=18"`
}

type UpdateUserRequest struct {
	Email    string `json:"email" validate:"omitempty,email"`
	Age sql.NullInt64 `json:"age" validate:"omitempty,gte=18"`
}

type DeleteUserRequest struct {
	ID int64 `json:"id" validate:"required"`
}
type GetUserRequest struct {
	ID int64 `json:"id" validate:"required,"`
}

type UserResponse struct {
	Message string `json:"msg"`
	Data []User `json:"data"`
}

