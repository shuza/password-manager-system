package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordMatch(t *testing.T) {
	t.Run("should match same password always", func(t *testing.T) {
		password := "some-password"
		passwordHash, err := HashPassword(password)
		assert.Nil(t, err)
		isMatch := CheckPasswordHash(password, passwordHash)
		assert.True(t, isMatch)
	})
}
