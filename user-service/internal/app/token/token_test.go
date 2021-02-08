package token

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"user-service/internal/app/model"
)

func TestService_Decode(t *testing.T) {
	t.Run("should decode the same user", func(t *testing.T) {
		user := model.User{
			ID:           1,
			Email:        "abcd@gmail.com",
			Password:     "123456",
			FullName:     "mr abcd",
			BusinessName: "my business",
		}

		s := NewService()
		token, err := s.Encode(user)
		assert.Nil(t, err)
		jwtClm, err := s.Decode(token)
		assert.Nil(t, err)
		assert.EqualValues(t, user, jwtClm.User)
	})

	t.Run("should return token parse error", func(t *testing.T) {
		token := ""
		s := NewService()
		jwtClm, err := s.Decode(token)
		assert.NotNil(t, err)
		assert.Nil(t, jwtClm)
	})
}
