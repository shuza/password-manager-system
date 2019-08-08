package service

import "github.com/stretchr/testify/mock"

type AuthMock struct {
	mock.Mock
}

func (m *AuthMock) GetUserId(token string) (int, error) {
	args := m.Mock.Called(token)
	if args.Get(1) != nil {
		return 0, args.Get(1).(error)
	}

	return args.Get(0).(int), nil
}
