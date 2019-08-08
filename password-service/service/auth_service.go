package service

type IAuthService interface {
	GetUserId(token string) (int, error)
}

var AuthService IAuthService
