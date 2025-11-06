package server

import (
	"net/http"
	"vadilatorgolang/internal/user" // Import package user
)

// NewRouter khởi tạo và trả về *http.ServeMux đã cấu hình
func NewRouter(userHandler *user.UserHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// Đăng ký route cho User
	// CÁC ROUTE KHÔNG CÓ ID
	mux.HandleFunc("POST /user", userHandler.CreateUserHandler)
	mux.HandleFunc("GET /user", userHandler.GetAllUserHandler)

	
	// 'GET /user/get/123'
	mux.HandleFunc("GET /user/{id}", userHandler.GetUserByIDHandler)

	// 'PUT /user/update/123'
	mux.HandleFunc("PUT /user/{id}", userHandler.UpdateUserHandler)

	// 'DELETE /user/delete/123'
	mux.HandleFunc("DELETE /user/{id}", userHandler.DeleteUserHandler)

	return mux
}

