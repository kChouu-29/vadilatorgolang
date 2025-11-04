package user

type UserController struct {
	Repo UserRepository
}

func NewUserController(r UserRepository) *UserController{
	return &UserController{Repo: r}
}

// Create 
func (u *UserController) CreateUser(user *User) error{
	return u.Repo.CreateUser(user)
}
// GetAllContact
func(u *UserController) GetAllContact() ([]User,error){
	return u.Repo.GetAllUser()
}
//GetByID
func(u *UserController) GetUserByID(id int64) (*User,error){
	return u.Repo.GetUserByID(id)
}
//Update
func(u *UserController) UpdateUserByID(user *User) error{
	return u.Repo.UpdateUserByID(user)
}
//DeleteByID
func (u *UserController) DeleteByID(id int64) error{
	return u.Repo.DeleteUserByID(id)
}