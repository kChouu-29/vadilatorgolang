package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"vadilatorgolang/package/logger"
	customValidator "vadilatorgolang/package/validator"

	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	Ctrl *UserController
}

func NewUserHandler(u *UserController) *UserHandler {
	return &UserHandler{Ctrl: u}
}

// ================== HANDLERS ===================

// CreateUserHandler
func (u *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	logger.TraceLogger.Printf("→ Bắt đầu CreateUserHandler. Request: %s %s", r.Method, r.URL.Path)

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.WarnLogger.Printf("Decode error: %v. Request: %s %s", err, r.Method, r.URL.Path)
		u.errorJson(w, http.StatusBadRequest, "Request body không hợp lệ")
		return
	}

	logger.DebugLogger.Printf("Body sau decode: %+v", req)

	// Validate
	if err := customValidator.ValidateStruct(req); err != nil {
		u.validationErrorJson(w, err, r)
		return
	}

	newUser := &User{
		UserName:  req.UserName,
		Email:     req.Email,
		Age:       req.Age,
		CreatedAt: time.Now(),
	}

	if err := u.Ctrl.CreateUser(newUser); err != nil {
		logger.ErrorLogger.Printf("Lỗi CreateUser: %v. Request: %s %s", err, r.Method, r.URL.Path)
		u.errorJson(w, http.StatusInternalServerError, "Không thể tạo user: "+err.Error())
		return
	}

	logger.InfoLogger.Printf("Tạo user thành công: %+v. Request: %s %s", newUser, r.Method, r.URL.Path)
	u.writeJson(w, http.StatusCreated, UserResponse{
		Message: "Tạo user thành công",
		Data:    []User{*newUser},
	})

	logger.TraceLogger.Printf("← Kết thúc CreateUserHandler. Request: %s %s", r.Method, r.URL.Path)
}

// GetUserByIDHandler
func (u *UserHandler) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	logger.TraceLogger.Printf("→ Bắt đầu GetUserByIDHandler. Request: %s %s", r.Method, r.URL.Path)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.WarnLogger.Printf("Invalid ID format: %s. Request: %s %s", idStr, r.Method, r.URL.Path)
		u.errorJson(w, http.StatusBadRequest, "Invalid ID format")
		return
	}
	if id <= 0 {
		logger.WarnLogger.Printf("ID must be positive: %d. Request: %s %s", id, r.Method, r.URL.Path)
		u.errorJson(w, http.StatusBadRequest, "ID must be a positive integer")
		return
	}
	user, err := u.Ctrl.GetUserByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.WarnLogger.Printf("Không tìm thấy user ID %d. Request: %s %s", id, r.Method, r.URL.Path)
			u.errorJson(w, http.StatusNotFound, "Không tìm thấy user")
		} else {
			logger.ErrorLogger.Printf("Lỗi GetUserByID %d: %v. Request: %s %s", id, err, r.Method, r.URL.Path)
			u.errorJson(w, http.StatusInternalServerError, "Lỗi truy vấn: "+err.Error())
		}
		return
	}

	logger.InfoLogger.Printf("Lấy user ID %d thành công. Request: %s %s", id, r.Method, r.URL.Path)
	u.writeJson(w, http.StatusOK, UserResponse{
		Message: "Lấy user thành công",
		Data:    []User{*user},
	})

	logger.TraceLogger.Printf("← Kết thúc GetUserByIDHandler. Request: %s %s", r.Method, r.URL.Path)
}

// GetAllUserHandler
func (u *UserHandler) GetAllUserHandler(w http.ResponseWriter, r *http.Request) {
	logger.TraceLogger.Printf("→ Bắt đầu GetAllUserHandler. Request: %s %s", r.Method, r.URL.Path)

	users, err := u.Ctrl.GetAllContact()
	if err != nil {
		logger.ErrorLogger.Printf("Lỗi GetAllContact: %v. Request: %s %s", err, r.Method, r.URL.Path)
		u.errorJson(w, http.StatusInternalServerError, "Lỗi lấy danh sách user: "+err.Error())
		return
	}

	logger.InfoLogger.Printf("Lấy tất cả user thành công (%d user). Request: %s %s", len(users), r.Method, r.URL.Path)
	u.writeJson(w, http.StatusOK, UserResponse{
		Message: "Lấy tất cả user thành công",
		Data:    users,
	})

	logger.TraceLogger.Printf("← Kết thúc GetAllUserHandler. Request: %s %s", r.Method, r.URL.Path)
}

// UpdateUserHandler
func (u *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	logger.TraceLogger.Printf("→ Bắt đầu UpdateUserHandler. Request: %s %s", r.Method, r.URL.Path)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.WarnLogger.Printf("Invalid ID format: %s. Request: %s %s", idStr, r.Method, r.URL.Path)
		u.errorJson(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.WarnLogger.Printf("Invalid JSON body: %v. Request: %s %s", err, r.Method, r.URL.Path)
		u.errorJson(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	user.ID = id

	if err := u.Ctrl.UpdateUserByID(&user); err != nil {
		if err == sql.ErrNoRows {
			logger.WarnLogger.Printf("Không tìm thấy user ID %d để cập nhật. Request: %s %s", id, r.Method, r.URL.Path)
			u.errorJson(w, http.StatusNotFound, err.Error())
		} else {
			logger.ErrorLogger.Printf("Lỗi UpdateUserByID %d: %v. Request: %s %s", id, err, r.Method, r.URL.Path)
			u.errorJson(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	logger.InfoLogger.Printf("Cập nhật user ID %d thành công. Request: %s %s", id, r.Method, r.URL.Path)
	u.writeJson(w, http.StatusOK, UserResponse{
		Message: "Update user successful",
		Data:    []User{user},
	})

	logger.TraceLogger.Printf("← Kết thúc UpdateUserHandler. Request: %s %s", r.Method, r.URL.Path)
}

// DeleteUserHandler
func (u *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	logger.TraceLogger.Printf("→ Bắt đầu DeleteUserHandler. Request: %s %s", r.Method, r.URL.Path)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.WarnLogger.Printf("Invalid ID format: %s. Request: %s %s", idStr, r.Method, r.URL.Path)
		u.errorJson(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	if err := u.Ctrl.DeleteByID(id); err != nil {
		if err == sql.ErrNoRows {
			logger.WarnLogger.Printf("Không tìm thấy user ID %d để xóa. Request: %s %s", id, r.Method, r.URL.Path)
			u.errorJson(w, http.StatusNotFound, "Không tìm thấy user để xóa")
		} else {
			logger.ErrorLogger.Printf("Lỗi xóa user ID %d: %v. Request: %s %s", id, err, r.Method, r.URL.Path)
			u.errorJson(w, http.StatusInternalServerError, "Lỗi xóa user: "+err.Error())
		}
		return
	}

	logger.InfoLogger.Printf("Xóa user ID %d thành công. Request: %s %s", id, r.Method, r.URL.Path)
	u.writeJson(w, http.StatusOK, UserResponse{
		Message: "Xóa user thành công",
		Data:    nil,
	})

	logger.TraceLogger.Printf("← Kết thúc DeleteUserHandler. Request: %s %s", r.Method, r.URL.Path)
}

// ================== HELPER FUNCTIONS ===================

func (u *UserHandler) writeJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (u *UserHandler) errorJson(w http.ResponseWriter, status int, message string) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (u *UserHandler) validationErrorJson(w http.ResponseWriter, err error, r *http.Request) {
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		msgs := make(map[string]string)
		for _, e := range ve {
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

		logger.WarnLogger.Printf("Lỗi Validation: %v. Request: %s %s", msgs, r.Method, r.URL.Path)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": msgs})
		return
	}

	logger.ErrorLogger.Printf("Lỗi Validation (unknown): %v. Request: %s %s", err, r.Method, r.URL.Path)
	u.errorJson(w, http.StatusBadRequest, "Data not correct")
}
