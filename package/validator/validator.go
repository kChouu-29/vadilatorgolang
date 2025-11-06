package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	validate      = validator.New()
	usernameRegex = regexp.MustCompile("^[a-zA-Z0-9_]+$")
)

// BỎ HÀM INIT() ĐI

// THÊM HÀM PUBLIC NÀY
// RegisterCustomValidations sẽ đăng ký tất cả các quy tắc của bạn
func RegisterCustomValidations() {
	err := validate.RegisterValidation("username_chars", validateUsernameChars)
	if err != nil {
		panic("Không thể đăng ký custom validation 'username_chars': " + err.Error())
	}
}

// Hàm custom (giữ nguyên, có thể để private 'v' thường)
func validateUsernameChars(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if username == "" {
		return true
	}
	return usernameRegex.MatchString(username)
}

// Hàm ValidateStruct (giữ nguyên)
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}
