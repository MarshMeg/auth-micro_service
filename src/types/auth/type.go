package auth

type Token struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
	TTL    int    `json:"ttl"`
}
