package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"net/http"
	"strconv" // Dùng để chuyển đổi ID
	"time"

	"vadilatorgolang/package/logger" // <-- IMPORT MỚI
	customValidator "vadilatorgolang/package/validator"

	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	Ctrl *UserController
}

func NewUserHandler(u *UserController) *UserHandler {
	return &UserHandler{Ctrl: u}
}

// === CÁC HANDLER ===

// CreateUserHandler
func (u *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// SỬA: Dùng ErrorLogger
		logger.ErrorLogger.Printf("Decode error: %v. Request: %s %s", err, r.Method, r.URL.Path)
		u.errorJson(w, http.StatusBadRequest, "Request body không hợp lệ")
		return
	}

	// Validate
	if err := customValidator.ValidateStruct(req); err != nil {
		u.validationErrorJson(w, err, r) // Truyền r vào
		return
	}

	newUser := &User{
		UserName:  req.UserName,
		Email:     req.Email,
		Age:       req.Age,
		CreatedAt: time.Now(),
	}

	if err := u.Ctrl.CreateUser(newUser); err != nil {
		// SỬA: Dùng ErrorLogger
		logger.ErrorLogger.Printf("Lỗi CreateUser: %v. Request: %s %s", err, r.Method, r.URL.Path)
		u.errorJson(w, http.StatusInternalServerError, "Không thể tạo user: "+err.Error())
		return
	}

	// Ghi log thành công
	logger.SuccessLogger.Printf("Tạo user thành công. Request: %s %s", r.Method, r.URL.Path)
	u.writeJson(w, http.StatusCreated, UserResponse{
		Message: "Tạo user thành công",
		Data:    []User{*newUser},
	})
}

// GetUserByIDHandler
func (u *UserHandler) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.ErrorLogger.Printf("Invalid ID format: %s. Request: %s %s", idStr, r.Method, r.URL.Path)
		u.errorJson(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	user, err := u.Ctrl.GetUserByID(id)
	if err != nil {
		if errors.Is(err, errors.New("sql: no rows in result set")) {
			logger.ErrorLogger.Printf("Không tìm thấy user ID %d. Request: %s %s", id, r.Method, r.URL.Path)
			u.errorJson(w, http.StatusNotFound, "Không tìm thấy user")
		} else {
			logger.ErrorLogger.Printf("Lỗi GetUserByID %d: %v. Request: %s %s", id, err, r.Method, r.URL.Path)
			u.errorJson(w, http.StatusInternalServerError, "Lỗi truy vấn: "+err.Error())
		}
		return
	}

	logger.SuccessLogger.Printf("Lấy user ID %d thành công. Request: %s %s", id, r.Method, r.URL.Path)
	u.writeJson(w, http.StatusOK, UserResponse{
		Message: "Lấy user thành công",
		Data:    []User{*user},
	})
}

// GetAllUserHandler
func (u *UserHandler) GetAllUserHandler(w http.ResponseWriter, r *http.Request) {
	users, err := u.Ctrl.GetAllContact()
	if err != nil {
		logger.ErrorLogger.Printf("Lỗi GetAllContact: %v. Request: %s %s", err, r.Method, r.URL.Path)
		u.errorJson(w, http.StatusInternalServerError, "Lỗi lấy danh sách user: "+err.Error())
		return
	}

	logger.SuccessLogger.Printf("Lấy tất cả user thành công. Request: %s %s", r.Method, r.URL.Path)
	u.writeJson(w, http.StatusOK, UserResponse{
		Message: "Lấy tất cả user thành công",
		Data:    users,
	})
}

// UpdateUserHandler
func (u *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		u.errorJson(w, http.StatusBadRequest, "Invalid ID format")
		return
	}
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		u.errorJson(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	user.ID = id
	if err := u.Ctrl.UpdateUserByID(&user); err != nil {
		if err == sql.ErrNoRows {
			logger.ErrorLogger.Printf("Lỗi UpdateContact: %v.Request: %s %s",err,r.Method,r.URL.Path)
			u.errorJson(w, http.StatusNotFound, err.Error())
		} else {
			u.errorJson(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
		logger.SuccessLogger.Printf("Cập nhật user thành công. Request: %s %s", r.Method, r.URL.Path)
	u.writeJson(w, http.StatusOK, UserResponse{
		Message: "Update user successful",
		Data:    []User{user},
	})
}

// DeleteUserHandler
func (u *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// (Tương tự, bạn hãy thêm logger.SuccessLogger và logger.ErrorLogger)
	// (Code của bạn ở đây)
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		u.errorJson(w, http.StatusBadRequest, "Invalid ID format")
		return
	}
	if err := u.Ctrl.DeleteByID(id); err != nil {
		if err == sql.ErrNoRows {
			logger.ErrorLogger.Printf("Không tìm thấy user ID %d. Request: %s %s", id, r.Method, r.URL.Path)
			u.errorJson(w, http.StatusNotFound, "Không tìm thấy user để xóa")
		} else {
			u.errorJson(w, http.StatusInternalServerError, "Lỗi xóa user: "+err.Error())
		}
		return
	}
	logger.SuccessLogger.Printf("Xóa user thành công. Request: %s %s", r.Method, r.URL.Path)
	u.writeJson(w, http.StatusOK, UserResponse{
		Message: "Xóa user thành công",
		Data:    nil,
	})
}

// === CÁC HÀM HELPER ===

// Sửa: Hàm này giờ chỉ để ghi JSON thành công
func (u *UserHandler) writeJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Sửa: Hàm này sẽ tự ghi log lỗi
func (u *UserHandler) errorJson(w http.ResponseWriter, status int, message string) {
	// (Chúng ta đã ghi log ở hàm gọi, nên ở đây chỉ cần write)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// Sửa: Hàm này sẽ tự ghi log validation
func (u *UserHandler) validationErrorJson(w http.ResponseWriter, err error, r *http.Request) {
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		msgs := make(map[string]string)
		for _, e := range ve {
			// (Giữ nguyên logic tạo message)
			switch e.Tag() {
			case "gte":
				msgs[e.Field()] = fmt.Sprintf("%s must be greater than or equal %s", e.Field(), e.Param())
			case "email":
				msgs[e.Field()] = fmt.Sprintf("%s not in the correct format", e.Field())
			case "required":
				msgs[e.Field()] = fmt.Sprintf("Trường '%s' là bắt buộc", e.Field())
			case "username_chars":
				msgs[e.Field()] = "Trường 'UserName' chỉ được chứa chữ cái, số và dấu gạch dưới"
			default:
				msgs[e.Field()] = fmt.Sprintf("Trường '%s' vi phạm quy tắc '%s'", e.Field(), e.Tag())
			}
		}

		// Ghi log chi tiết lỗi validation
		logger.ErrorLogger.Printf("Lỗi Validation: %v. Request: %s %s", msgs, r.Method, r.URL.Path)

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": msgs})
		return
	}

	// Lỗi validation không xác định
	logger.ErrorLogger.Printf("Lỗi Validation (unknown): %v. Request: %s %s", err, r.Method, r.URL.Path)
	u.errorJson(w, http.StatusBadRequest, "Data not correct")
}
