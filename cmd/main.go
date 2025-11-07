package main

import (
	
	"net/http"

	"vadilatorgolang/internal/user"
	"vadilatorgolang/package/database"
	"vadilatorgolang/package/logger"
	"vadilatorgolang/package/server"
	customValidator "vadilatorgolang/package/validator"
)

func main() {
	// 0. KHỞI TẠO LOGGER (Thêm bước này đầu tiên)
	// (Đảm bảo thư mục 'logs' đã được tạo)
	logger.InitLoggers()

	// 1. Kết nối DB
	db, err := database.ConnectDb()
	if err != nil {
		// Dùng ErrorLogger cho lỗi database
		logger.ErrorLogger.Fatal("Không thể kết nối tới database: ", err)
	}
	defer db.Close()
	// Ghi log thành công
	logger.SuccessLogger.Println("Kết nối database thành công.")

	customValidator.RegisterCustomValidations()
	logger.SuccessLogger.Println("Đã đăng ký custom validators.")

	// 3. Khởi tạo (Dependency Injection)
	userRepo := user.NewUserRepo(db)
	userCtrl := user.NewUserController(userRepo)
	userHandler := user.NewUserHandler(userCtrl)

	// 4. Khởi tạo Router
	router := server.NewRouter(userHandler)

	// 5. Khởi động server
	port := ":8080"
	logger.SuccessLogger.Printf("Server đang chạy tại http://localhost%s", port)

	err = http.ListenAndServe(port, router)
	if err != nil {
		logger.ErrorLogger.Fatal("Lỗi khi khởi động server: ", err)
	}
}
