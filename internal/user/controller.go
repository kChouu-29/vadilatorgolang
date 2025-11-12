package user

import (
	"database/sql"
	"errors"
)

// Controller giữ Repo (như file gốc của bạn)
type UserController struct {
	Repo UserRepository
}

// NewUserController nhận vào Repo (như file gốc của bạn)
func NewUserController(r UserRepository) *UserController {
	return &UserController{Repo: r}
}

// --- LOGIC NGHIỆP VỤ ĐƯỢC ĐẶT TRỰC TIẾP TẠI ĐÂY ---

// Create
func (u *UserController) CreateUser(user *User) error {
	// --- KIỂM TRA USERNAME ---
	existingUser, err := u.Repo.GetUserByUsername(user.UserName)
	if err != nil && err != sql.ErrNoRows {
		return err // Lỗi database
	}
	if existingUser != nil {
		return errors.New("username đã tồn tại") // Lỗi nghiệp vụ
	}

	// --- KIỂM TRA EMAIL ---
	existingUser, err = u.Repo.GetUserByEmail(user.Email)
	if err != nil && err != sql.ErrNoRows {
		return err // Lỗi database
	}
	if existingUser != nil {
		return errors.New("email đã tồn tại") // Lỗi nghiệp vụ
	}

	// Nếu mọi thứ ổn, gọi Repo
	return u.Repo.CreateUser(user)
}

// GetAllContact
func (u *UserController) GetAllContact() ([]User, error) {
	return u.Repo.GetAllUser()
}

// GetByID
func (u *UserController) GetUserByID(id int) (*User, error) {
	return u.Repo.GetUserByID(id)
}

// Update
func (u *UserController) UpdateUserByID(user *User) error {
	
	return u.Repo.UpdateUserByID(user)
}

// DeleteByID
func (u *UserController) DeleteByID(id int) error {
	return u.Repo.DeleteUserByID(id)
}
