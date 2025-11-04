package user

import (
	"database/sql"
	"fmt"
)

type UserRepository interface {
	CreateUser(c *User) error
	GetUserByID(id int64) (*User, error)
	GetAllUser() ([]User, error)
	UpdateUserByID(u *User) error
	DeleteUserByID(id int64) error
}

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepository {
	return &UserRepo{DB: db}
}

// Create
func (r *UserRepo) CreateUser(c *User) error {
	_, err := r.DB.Exec("insert into nguoi_dung(usename,email,age,created_at) values(?,?,?)", c.UserName, c.Email, c.Age, c.CreatedAt)
	if err != nil {
		return err
	}
	return err
}

//Get by ID

func (r *UserRepo) GetUserByID(id int64) (*User, error) {
	row := r.DB.QueryRow("select id,username,email,age,created_at from nguoi_dung where id=?", id)
	var c User
	if err := row.Scan(&c.ID, &c.UserName, &c.Email, &c.Age, &c.CreatedAt); err != nil {
		return nil, err
	}
	return &c, nil
}

// Get all
func (r *UserRepo) GetAllUser() ([]User, error) {
	row, err := r.DB.Query("select id,username,email,age,created_at from nguoi_dung")
	if err != nil {
		return nil, err
	}
	var c []User
	for row.Next() {
		var p User
		if err := row.Scan(&p.ID, &p.UserName, &p.Email, &p.Age, &p.CreatedAt); err != nil {
			return nil, err
		}
	}
	return c, nil
}

//Update By ID

func (r *UserRepo) UpdateUserByID(c *User) error {
	res, err := r.DB.Exec("update nguoi_dung set username=?,email=?,age=?,created_at=?where id=?", c.UserName, c.Email, c.Age, c.CreatedAt, c.ID)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return fmt.Errorf("contact with ID %d not found", c.ID)
	}
	return nil
}

//Delete By Id

func (r *UserRepo) DeleteUserByID(id int64) error {
	res, err := r.DB.Exec("Delete from nguoi_dung where id=?", id)
	if err != nil {
		return nil
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("contact with ID %d not found",id)
	}
	return nil
}
