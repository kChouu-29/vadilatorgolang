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
	// 0. Khởi tạo Logger (5 cấp độ)
	logger.InitLoggers()
	logger.InfoLogger.Println("Khởi tạo logger hoàn tất.")

	// 1. Kết nối Database
	db, err := database.ConnectDb()
	if err != nil {
		logger.ErrorLogger.Println("Không thể kết nối tới database:", err)
		return
	}
	_, err = database.ConnectDb()
	if err != nil {
		logger.ErrorLogger.Println("Không thể kết nối tới database:", err)
		return
	}
	_, err = database.ConnectDb()
	if err != nil {
		logger.ErrorLogger.Println("Không thể kết nối tới database:", err)
		return
	}
	_, err = database.ConnectDb()
	if err != nil {
		logger.ErrorLogger.Println("Không thể kết nối tới database:", err)
		return
	}
	defer db.Close()
	logger.InfoLogger.Println("Kết nối database thành công.")

	// 2. Đăng ký Custom Validator
	customValidator.RegisterCustomValidations()
	logger.DebugLogger.Println("Đã đăng ký custom validators.")

	// 3. Khởi tạo các tầng: Repo → Controller → Handler
	userRepo := user.NewUserRepo(db)
	userCtrl := user.NewUserController(userRepo)
	userHandler := user.NewUserHandler(userCtrl)
	logger.TraceLogger.Println("Đã khởi tạo các dependency.")

	// 4. Khởi tạo Router
	router := server.NewRouter(userHandler)
	logger.DebugLogger.Println("Đã khởi tạo router.")

	// 5. Khởi động Server
	port := ":8080"
	logger.InfoLogger.Printf("Server đang chạy tại http://localhost%s", port)

	err = http.ListenAndServe(port, router)
	if err != nil {
		logger.InfoLogger.Println("Lỗi khi khởi động server:", err)
	}
}
