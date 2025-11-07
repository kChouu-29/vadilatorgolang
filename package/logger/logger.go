package logger

import (
	"io"
	"log"
	"os"
)

var (

	SuccessLogger *log.Logger

	ErrorLogger *log.Logger
)


func InitLoggers() {
	
	successFile, err := os.OpenFile("log/success.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Không thể mở file log/success.log: ", err)
	}

	// 2. Tạo MultiWriter để ghi ra File VÀ Console (os.Stdout)
	successWriter := io.MultiWriter(successFile, os.Stdout)

	
	SuccessLogger = log.New(successWriter, "SUCCESS: ", log.Ldate|log.Ltime|log.Lshortfile)

	
	errorFile, err := os.OpenFile("log/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Không thể mở file log/error.log: ", err)
	}

	
	errorWriter := io.MultiWriter(errorFile, os.Stderr)

	
	ErrorLogger = log.New(errorWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	SuccessLogger.Println("Logger Thành Công đã khởi tạo.")
	ErrorLogger.Println("Logger Lỗi đã khởi tạo.")
}