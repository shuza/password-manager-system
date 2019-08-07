package db

import (
	"github.com/stretchr/testify/mock"
	"password-service/model"
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

func (m *MockDb) Delete(model interface{}) error {
	args := m.Called(model)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func (m *MockDb) GetByUserId(userId uint) ([]model.Password, error) {
	args := m.Mock.Called(userId)
	if args.Get(1) != nil {
		return []model.Password{}, args.Get(1).(error)
	}

	return args.Get(0).([]model.Password), nil
}

func (m *MockDb) GetById(id int, userId int) (model.Password, error) {
	args := m.Called(id, userId)
	if args.Get(1) != nil {
		return model.Password{}, args.Get(1).(error)
	}

	return args.Get(0).(model.Password), nil
}

func (m *MockDb) GetByUsernameAndPassword(username string, password string) (model.Password, error) {
	args := m.Called(username, password)
	if args.Get(1) != nil {
		return model.Password{}, args.Get(1).(error)
	}

	return args.Get(0).(model.Password), nil
}

func (m *MockDb) Close() {
}
