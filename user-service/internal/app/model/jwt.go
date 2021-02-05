package model

import (
	"github.com/dgrijalva/jwt-go"
)

type JwtCustomClaims struct {
	User User
	jwt.StandardClaims
}
