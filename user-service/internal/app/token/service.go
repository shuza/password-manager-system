package token

import (
	"fmt"
	"time"
	"user-service/internal/app/model"

	"github.com/dgrijalva/jwt-go"
)

var (
	//	Define a secure key string used
	//	as salt when hashing our tokens
	key = []byte("hashingpasswordismandatory")
)

type service struct{}

func NewService() *service {
	return &service{}
}

func (s *service) Encode(user model.User) (string, error) {
	expireToken := time.Now().Add(72 * time.Hour).Unix()

	//	Create claim
	claim := model.JwtCustomClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "auth.user",
		},
	}

	//	Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	//	Sign token and return
	return token.SignedString(key)
}

func (s *service) Decode(tokenStr string) (*model.JwtCustomClaims, error) {
	//	Parse token
	token, err := jwt.ParseWithClaims(tokenStr, &model.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	//	Validate the token and return custom claim
	if claim, ok := token.Claims.(*model.JwtCustomClaims); ok && token.Valid {
		return claim, nil
	}

	return nil, fmt.Errorf("token not found :%w", model.ErrNotFound)
}
