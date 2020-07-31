package model

type UserDetails struct {
	// 用户标识
	UserID   int64
	// 用户名
	Username string
	// 用户密码
	Password string
	// 用户具备的权限
	Authorities  []string
}

// IsMatch 用户是否匹配
func (userDetails *UserDetails) IsMatch(username, password string) bool {
	return userDetails.Username == username && userDetails.Password == password
}
