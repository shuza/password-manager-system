package db

import (
	"github.com/stretchr/testify/mock"
	"user-service/model"
)

type MockDb struct {
	mock.Mock
}

func (m *MockDb) Init() error {
	return nil
}

func (m *MockDb) Save(model interface{}) error {
	args := m.Called(model)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func (m *MockDb) GetByEmail(email string) (model.User, error) {
	args := m.Called(email)
	if args.Get(1) != nil {
		return args.Get(0).(model.User), args.Get(1).(error)
	}
	return args.Get(0).(model.User), nil
}

func (m *MockDb) Close() {
}
