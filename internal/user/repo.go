package user

import (
	"database/sql"
	"fmt"
)

// UserRepository là interface định nghĩa các phương thức cho database
type UserRepository interface {
	CreateUser(c *User) error
	GetUserByID(id int) (*User, error)
	GetAllUser() ([]User, error)
	UpdateUserByID(u *User) error
	DeleteUserByID(id int) error
	GetUserByEmail(email string) (*User, error)
	GetUserByUsername(username string) (*User, error)
}

// UserRepo là struct triển khai UserRepository
type UserRepo struct {
	DB *sql.DB
}

// NewUserRepo tạo một repository mới
func NewUserRepo(db *sql.DB) UserRepository {
	return &UserRepo{DB: db}
}

// Create
func (r *UserRepo) CreateUser(c *User) error {
	_, err := r.DB.Exec("insert into nguoi_dung(username,email,age,created_at) values(?,?,?,?)", c.UserName, c.Email, c.Age, c.CreatedAt)
	if err != nil {
		return err
	}
	return nil // Sửa: trả về nil khi thành công
}

// Get by ID
func (r *UserRepo) GetUserByID(id int) (*User, error) {
	row := r.DB.QueryRow("select id,username,email,age,created_at from nguoi_dung where id=?", id)
	var c User
	if err := row.Scan(&c.ID, &c.UserName, &c.Email, &c.Age, &c.CreatedAt); err != nil {
		return nil, err
	}
	return &c, nil
}

// === THÊM MỚI: Get by Email ===
// GetUserByEmail tìm người dùng bằng email
func (r *UserRepo) GetUserByEmail(email string) (*User, error) {
	row := r.DB.QueryRow("select id,username,email,age,created_at from nguoi_dung where email=?", email)
	var c User
	if err := row.Scan(&c.ID, &c.UserName, &c.Email, &c.Age, &c.CreatedAt); err != nil {
		// err ở đây có thể là 'sql.ErrNoRows' (không tìm thấy)
		return nil, err
	}
	return &c, nil // Tìm thấy user
}

// === KẾT THÚC THÊM MỚI ===

// Get all
func (r *UserRepo) GetAllUser() ([]User, error) {
	row, err := r.DB.Query("select id,username,email,age,created_at from nguoi_dung")
	if err != nil {
		return nil, err
	}
	defer row.Close() // Thêm: Đóng rows sau khi dùng

	var c []User
	for row.Next() {
		var p User
		if err := row.Scan(&p.ID, &p.UserName, &p.Email, &p.Age, &p.CreatedAt); err != nil {
			return nil, err
		}
		c = append(c, p) // Thêm: Phải append vào slice
	}
	return c, nil
}

// Update By ID
func (r *UserRepo) UpdateUserByID(c *User) error {
	res, err := r.DB.Exec("update nguoi_dung set username=?,email=?,age=?,created_at=? where id=?", c.UserName, c.Email, c.Age, c.CreatedAt, c.ID) // Sửa: Thêm khoảng trắng trước 'where'
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return fmt.Errorf("contact with ID %d not found", c.ID)
	}
	return nil
}

// Delete By Id
func (r *UserRepo) DeleteUserByID(id int) error {
	res, err := r.DB.Exec("Delete from nguoi_dung where id=?", id)
	if err != nil {
		return err // Sửa: Trả về err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("contact with ID %d not found", id)
	}
	return nil
}
func (r *UserRepo) GetUserByUsername(username string) (*User, error) {
	row := r.DB.QueryRow("select id,username,email,age,created_at from nguoi_dung where username=?", username)
	var c User
	if err := row.Scan(&c.ID, &c.UserName, &c.Email, &c.Age, &c.CreatedAt); err != nil {
		return nil, err // Trả về lỗi (ví dụ: sql.ErrNoRows)
	}
	return &c, nil // Tìm thấy user
}
