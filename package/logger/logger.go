package logger

import (
	"io"
	"log"
	"os"
)

var (
	TraceLogger *log.Logger
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
)

// InitLoggers khởi tạo các logger theo từng mức độ
func InitLoggers() {
	// Tạo thư mục log nếu chưa có
	if _, err := os.Stat("log"); os.IsNotExist(err) {
		os.Mkdir("log", 0755)
	}

	// TRACE 
	traceFile, err := os.OpenFile("log/trace.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Không thể mở file log/trace.log:", err)
	}
	traceWriter := io.MultiWriter(traceFile, os.Stdout)
	TraceLogger = log.New(traceWriter, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)

	// DEBUG 
	debugFile, err := os.OpenFile("log/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Không thể mở file log/debug.log:", err)
	}
	debugWriter := io.MultiWriter(debugFile, os.Stdout)
	DebugLogger = log.New(debugWriter, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

	// INFO 
	infoFile, err := os.OpenFile("log/info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Không thể mở file log/info.log:", err)
	}
	infoWriter := io.MultiWriter(infoFile, os.Stdout)
	InfoLogger = log.New(infoWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	// WARN 
	warnFile, err := os.OpenFile("log/warn.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Không thể mở file log/warn.log:", err)
	}
	warnWriter := io.MultiWriter(warnFile, os.Stdout)
	WarnLogger = log.New(warnWriter, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)

	// ERROR 
	errorFile, err := os.OpenFile("log/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Không thể mở file log/error.log:", err)
	}
	errorWriter := io.MultiWriter(errorFile, os.Stderr)
	ErrorLogger = log.New(errorWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Ghi log khi khởi tạo thành công
	InfoLogger.Println("Logger hệ thống đã được khởi tạo thành công.")
}
