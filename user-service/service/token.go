package service

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"user-service/model"
)

var (
	//	Define a secure key string used
	//	as salt when hashing our tokens
	key = []byte("hashingpasswordismandatory")
)

type CustomClaims struct {
	User model.User
	jwt.StandardClaims
}

type Authable interface {
	Decode(tokenStr string) (*CustomClaims, error)
	Encode(user model.User) (string, error)
}

type TokenService struct{}

func (s *TokenService) Encode(user model.User) (string, error) {
	expireToken := time.Now().Add(72 * time.Hour).Unix()

	//	Create claim
	claim := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "xendit.user",
		},
	}

	//	Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	//	Sign token and return
	return token.SignedString(key)
}

func (s *TokenService) Decode(tokenStr string) (*CustomClaims, error) {
	//	Parse token
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	//	Validate the token and return custom claim
	if claim, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claim, nil
	}

	return nil, err
}
