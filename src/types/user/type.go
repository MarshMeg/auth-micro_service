package user

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	RoleName string `json:"user_role"`
}

func (u *User) Return() *User {
	u.Password = "<secret>"
	return u
}
