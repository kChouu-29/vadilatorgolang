package main

import (
	"log"
	"net/http"

	"vadilatorgolang/internal/user"
	"vadilatorgolang/package/database"
	"vadilatorgolang/package/server"
	customValidator "vadilatorgolang/package/validator"
)

func main() {
	// 1. Kết nối DB
	db, err := database.ConnectDb()
	if err != nil {
		log.Fatal("Không thể kết nối tới database: ", err)
	}
	defer db.Close()

	customValidator.RegisterCustomValidations()
	log.Println("Đã đăng ký custom validators.")

	// 3. Khởi tạo (Dependency Injection)
	userRepo := user.NewUserRepo(db)
	userCtrl := user.NewUserController(userRepo)
	userHandler := user.NewUserHandler(userCtrl)

	// 4. Khởi tạo Router
	router := server.NewRouter(userHandler)

	// 5. Khởi động server
	port := ":8080"
	log.Printf("Server đang chạy tại http://localhost%s", port)

	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal("Lỗi khi khởi động server: ", err)
	}
}
