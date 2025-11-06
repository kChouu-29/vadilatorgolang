package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv" // Dùng để chuyển đổi ID
	"time"

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

// CreateUserHandler (Giữ nguyên, không dùng ID)
func (u *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Decode error: %v", err)
		u.errorJson(w, http.StatusBadRequest, "Request body không hợp lệ")
		return
	}

	// Validate
	if err := customValidator.ValidateStruct(req); err != nil {
		u.validationErrorJson(w, err)
		return
	}

	newUser := &User{
		UserName:  req.UserName,
		Email:     req.Email,
		Age:       req.Age,
		CreatedAt: time.Now(),
	}

	if err := u.Ctrl.CreateUser(newUser); err != nil {
		u.errorJson(w, http.StatusInternalServerError, "Không thể tạo user: "+err.Error())
		return
	}

	u.writeJson(w, http.StatusCreated, UserResponse{
		Message: "Tạo user thành công",
		Data:    []User{*newUser},
	})
}

// GetUserByIDHandler (Phong cách mới)
func (u *UserHandler) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Lấy ID từ URL path (e.g., /user/get/1)
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		u.errorJson(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	// 2. Gọi Controller
	user, err := u.Ctrl.GetUserByID(id)
	if err != nil {
		if errors.Is(err, errors.New("sql: no rows in result set")) {
			u.errorJson(w, http.StatusNotFound, "Không tìm thấy user")
		} else {
			u.errorJson(w, http.StatusInternalServerError, "Lỗi truy vấn: "+err.Error())
		}
		return
	}

	// 3. Trả về thành công
	u.writeJson(w, http.StatusOK, UserResponse{
		Message: "Lấy user thành công",
		Data:    []User{*user},
	})
}

// GetAllUserHandler (Giữ nguyên, không dùng ID)
func (u *UserHandler) GetAllUserHandler(w http.ResponseWriter, r *http.Request) {
	users, err := u.Ctrl.GetAllContact()
	if err != nil {
		u.errorJson(w, http.StatusInternalServerError, "Lỗi lấy danh sách user: "+err.Error())
		return
	}

	u.writeJson(w, http.StatusOK, UserResponse{
		Message: "Lấy tất cả user thành công",
		Data:    users,
	})
}

// UpdateUserHandler (Phong cách mới giống hệt ví dụ của bạn)
func (u *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Lấy ID từ URL path (e.g., /user/update/1)
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr) // Dùng ParseInt vì ID là int64
	if err != nil {
		u.errorJson(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	// 2. Decode JSON body vào struct User
	// (Lưu ý: Bỏ qua DTO 'UpdateUserRequest' để giống hệt ví dụ của bạn)
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		u.errorJson(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	// 3. Gán ID từ URL (quan trọng)
	user.ID = id

	// 4. Gọi Controller
	// (Lưu ý: Repo của bạn cập nhật tất cả các trường)
	if err := u.Ctrl.UpdateUserByID(&user); err != nil {
		// Kiểm tra lỗi 'not found' từ repo
		if err.Error() == fmt.Sprintf("contact with ID %d not found", user.ID) {
			u.errorJson(w, http.StatusNotFound, err.Error())
		} else {
			u.errorJson(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// 5. Trả về thành công
	u.writeJson(w, http.StatusOK, UserResponse{
		Message: "Update user successful",
		Data:    []User{user},
	})
}

// DeleteUserHandler (Phong cách mới)
func (u *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Lấy ID từ URL path (e.g., /user/delete/1)
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr) // Dùng ParseInt vì ID là int64
	if err != nil {
		u.errorJson(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	// 2. Gọi Controller
	if err := u.Ctrl.DeleteByID(id); err != nil {
		if err.Error() == fmt.Sprintf("contact with ID %d not found", id) {
			u.errorJson(w, http.StatusNotFound, "Không tìm thấy user để xóa")
		} else {
			u.errorJson(w, http.StatusInternalServerError, "Lỗi xóa user: "+err.Error())
		}
		return
	}

	// 3. Trả về thành công
	u.writeJson(w, http.StatusOK, UserResponse{
		Message: "Xóa user thành công",
		Data:    nil,
	})
}

// === CÁC HÀM HELPER ===

func (u *UserHandler) writeJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (u *UserHandler) errorJson(w http.ResponseWriter, status int, message string) {
	u.writeJson(w, status, map[string]string{"error": message})
}

func (u *UserHandler) validationErrorJson(w http.ResponseWriter, err error) {
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
		u.writeJson(w, http.StatusBadRequest, map[string]interface{}{
			"error": msgs,
		})
		return
	}
	u.errorJson(w, http.StatusBadRequest, "Data not correct")
}
