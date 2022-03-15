package user

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"user-service/internal/app/model"
	"user-service/internal/app/service/mocks"
	"user-service/internal/app/util"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		desc     string
		payload  model.User
		mockRepo func() *mocks.MockUserRepository
		expErr   error
	}{
		{
			desc: "should return success",
			payload: model.User{
				Email:        "abcd@gmail.com",
				Password:     "123456",
				FullName:     "mr abcd",
				BusinessName: "business-1",
			},
			mockRepo: func() *mocks.MockUserRepository {
				r := mocks.NewMockUserRepository(ctrl)
				r.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(nil)
				return r
			},
			expErr: nil,
		},
		{
			desc:    "should return invalid user request error",
			payload: model.User{},
			mockRepo: func() *mocks.MockUserRepository {
				return mocks.NewMockUserRepository(ctrl)
			},
			expErr: fmt.Errorf("invalid user request :%w", model.ErrInvalid),
		},
		{
			desc: "should return db error",
			payload: model.User{
				Email:        "email",
				Password:     "12345",
				FullName:     "asdf",
				BusinessName: "asdf",
			},
			mockRepo: func() *mocks.MockUserRepository {
				r := mocks.NewMockUserRepository(ctrl)
				r.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(errors.New("db-error"))
				return r
			},
			expErr: errors.New("db-error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			s := NewService(tc.mockRepo())
			err := s.CreateUser(context.Background(), tc.payload)
			assert.Equal(t, tc.expErr, err)
		})
	}
}

func TestService_GetUserByEmailAndPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	password, _ := util.HashPassword("123456")
	user := model.User{
		ID:       1,
		Email:    "abcd@gmail.com",
		Password: password,
	}

	testCases := []struct {
		desc     string
		email    string
		password string
		mockRepo func() *mocks.MockUserRepository
		expErr   error
		expUser  model.User
	}{
		{
			desc:     "should return success",
			email:    "abcd@gmail.com",
			password: "123456",
			mockRepo: func() *mocks.MockUserRepository {
				r := mocks.NewMockUserRepository(ctrl)
				r.EXPECT().GetUserByEmail(gomock.Any(), "abcd@gmail.com").Return(user, nil)
				return r
			},
			expErr:  nil,
			expUser: user,
		},
		{
			desc:     "should return invalid request error",
			email:    "",
			password: "",
			mockRepo: func() *mocks.MockUserRepository {
				return mocks.NewMockUserRepository(ctrl)
			},
			expErr:  fmt.Errorf("invalid login request :%w", model.ErrInvalid),
			expUser: model.User{},
		},
		{
			desc:     "should return DB error",
			email:    "abcd@gmail.com",
			password: "123456",
			mockRepo: func() *mocks.MockUserRepository {
				r := mocks.NewMockUserRepository(ctrl)
				r.EXPECT().GetUserByEmail(gomock.Any(), "abcd@gmail.com").Return(model.User{}, errors.New("db-error"))
				return r
			},
			expErr:  errors.New("db-error"),
			expUser: model.User{},
		},
		{
			desc:     "should return wrong password error",
			email:    "abcd@gmail.com",
			password: "wrong-password",
			mockRepo: func() *mocks.MockUserRepository {
				r := mocks.NewMockUserRepository(ctrl)
				r.EXPECT().GetUserByEmail(gomock.Any(), "abcd@gmail.com").Return(user, nil)
				return r
			},
			expErr:  fmt.Errorf("wrong password :%w", model.ErrInvalid),
			expUser: model.User{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			s := NewService(tc.mockRepo())
			user, err := s.GetUserByEmailAndPassword(context.Background(), tc.email, tc.password)
			assert.Equal(t, tc.expErr, err)
			assert.EqualValues(t, tc.expUser, user)
		})
	}
}
