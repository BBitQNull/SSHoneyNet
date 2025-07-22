package auth

type AuthService interface {
	PasswordValidator(username, password string) bool
}

type SimpleAuthService struct{}

func (a SimpleAuthService) PasswordValidator(username, password string) bool {
	return username == "root" && password == "123456"
}

// 可更改认证逻辑，如查询数据库
