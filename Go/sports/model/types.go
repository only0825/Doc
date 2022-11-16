package model

// 结构体定义

type User struct {
	Id          int    `json:"id" db:"id"`
	Username    string `json:"username" db:"username"`
	Password    string `json:"password" db:"password"`
	State       int    `json:"state" db:"state"`
	LastLoginIp string `json:"last_login_ip" db:"last_login_ip"`
	CreatedAt   int64  `json:"created_at" db:"created_at"`
	UpdatedAt   int64  `json:"updated_at" db:"updated_at"`
}

type UserLoginResp struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Domain   string `json:"domain"`
}
