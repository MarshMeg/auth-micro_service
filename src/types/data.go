package types

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     int    `json:"user_role"`
}

type Token struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
	TTL    int    `json:"ttl"`
}
