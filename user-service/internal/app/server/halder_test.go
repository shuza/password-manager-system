package server

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"user-service/internal/app/model"
	"user-service/internal/app/service/mocks"
)

func TestSignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		desc          string
		payload       string
		mockSvc       func() *mocks.MockUserService
		expStatusCode int
	}{
		{
			desc:    "should success",
			payload: `{ "email": "abcd@gmail.com", "full_name": "Mr. abcd", "password": "123456", "business_name": "my business 1" }`,
			mockSvc: func() *mocks.MockUserService {
				s := mocks.NewMockUserService(ctrl)
				s.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil)
				return s
			},
			expStatusCode: http.StatusCreated,
		},
		{
			desc:    "should return decode error",
			payload: `------------`,
			mockSvc: func() *mocks.MockUserService {
				return mocks.NewMockUserService(ctrl)
			},
			expStatusCode: http.StatusUnprocessableEntity,
		},
		{
			desc:    "should return invalid user error",
			payload: `{ "email": "abcd@gmail.com", "full_name": "Mr. abcd", "password": "123456", "business_name": "my business 1" }`,
			mockSvc: func() *mocks.MockUserService {
				s := mocks.NewMockUserService(ctrl)
				s.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(model.ErrInvalid)
				return s
			},
			expStatusCode: http.StatusBadRequest,
		},
		{
			desc:    "should return internal server error",
			payload: `{ "email": "abcd@gmail.com", "full_name": "Mr. abcd", "password": "123456", "business_name": "my business 1" }`,
			mockSvc: func() *mocks.MockUserService {
				s := mocks.NewMockUserService(ctrl)
				s.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(errors.New("server-error"))
				return s
			},
			expStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			s := NewServer(":8080", tc.mockSvc(), nil)

			w := httptest.NewRecorder()
			body := strings.NewReader(tc.payload)
			r := httptest.NewRequest(http.MethodPost, "/api/v1/signup", body)

			router := mux.NewRouter()
			router.Methods(http.MethodPost).Path("/api/v1/signup").HandlerFunc(s.signUp)
			router.ServeHTTP(w, r)
			assert.Equal(t, tc.expStatusCode, w.Code)
		})
	}
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := model.User{
		ID:           1,
		Email:        "abcd@gmail.com",
		Password:     "123456",
		FullName:     "Mr abcd",
		BusinessName: "my business",
	}

	testCases := []struct {
		desc          string
		payload       string
		mockUserSvc   func() *mocks.MockUserService
		mockAuthSvc   func() *mocks.MockAuthService
		expStatusCode int
	}{
		{
			desc:    "should success",
			payload: `{ "email": "abcd@gmail.com", "password": "123456" }`,
			mockUserSvc: func() *mocks.MockUserService {
				s := mocks.NewMockUserService(ctrl)
				s.EXPECT().GetUserByEmailAndPassword(gomock.Any(), "abcd@gmail.com", "123456").Return(user, nil)
				return s
			},
			mockAuthSvc: func() *mocks.MockAuthService {
				s := mocks.NewMockAuthService(ctrl)
				s.EXPECT().Encode(user).Return("auto-token", nil)
				return s
			},
			expStatusCode: http.StatusOK,
		},
		{
			desc:    "should return decode error",
			payload: `------------`,
			mockUserSvc: func() *mocks.MockUserService {
				return mocks.NewMockUserService(ctrl)
			},
			mockAuthSvc: func() *mocks.MockAuthService {
				return mocks.NewMockAuthService(ctrl)
			},
			expStatusCode: http.StatusUnprocessableEntity,
		},
		{
			desc:    "should return invalid credentials error",
			payload: `{ "email": "abcd@gmail.com", "password": "123456" }`,
			mockUserSvc: func() *mocks.MockUserService {
				s := mocks.NewMockUserService(ctrl)
				s.EXPECT().GetUserByEmailAndPassword(gomock.Any(), "abcd@gmail.com", "123456").Return(model.User{}, model.ErrInvalid)
				return s
			},
			mockAuthSvc: func() *mocks.MockAuthService {
				return mocks.NewMockAuthService(ctrl)
			},
			expStatusCode: http.StatusBadRequest,
		},
		{
			desc:    "should return internal server error",
			payload: `{ "email": "abcd@gmail.com", "password": "123456" }`,
			mockUserSvc: func() *mocks.MockUserService {
				s := mocks.NewMockUserService(ctrl)
				s.EXPECT().GetUserByEmailAndPassword(gomock.Any(), "abcd@gmail.com", "123456").Return(model.User{}, errors.New("server-error"))
				return s
			},
			mockAuthSvc: func() *mocks.MockAuthService {
				return mocks.NewMockAuthService(ctrl)
			},
			expStatusCode: http.StatusInternalServerError,
		},
		{
			desc:    "should return internal server error for jwt",
			payload: `{ "email": "abcd@gmail.com", "password": "123456" }`,
			mockUserSvc: func() *mocks.MockUserService {
				s := mocks.NewMockUserService(ctrl)
				s.EXPECT().GetUserByEmailAndPassword(gomock.Any(), "abcd@gmail.com", "123456").Return(user, nil)
				return s
			},
			mockAuthSvc: func() *mocks.MockAuthService {
				s := mocks.NewMockAuthService(ctrl)
				s.EXPECT().Encode(user).Return("", errors.New("jwt-error"))
				return s
			},
			expStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			s := NewServer(":8080", tc.mockUserSvc(), tc.mockAuthSvc())

			w := httptest.NewRecorder()
			body := strings.NewReader(tc.payload)
			r := httptest.NewRequest(http.MethodPost, "/api/v1/login", body)

			router := mux.NewRouter()
			router.Methods(http.MethodPost).Path("/api/v1/login").HandlerFunc(s.login)
			router.ServeHTTP(w, r)
			assert.Equal(t, tc.expStatusCode, w.Code)
		})
	}
}

func TestTokenValidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	claim := &model.JwtCustomClaims{
		User: model.User{
			ID:           1,
			Email:        "abcd@gmail.com",
			Password:     "123456",
			FullName:     "mr abcd",
			BusinessName: "my business",
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			Issuer:    "auth.user",
		},
	}

	testCases := []struct {
		desc          string
		token         string
		mockSvc       func() *mocks.MockAuthService
		expStatusCode int
	}{
		{
			desc:  "should success",
			token: "valid-token",
			mockSvc: func() *mocks.MockAuthService {
				s := mocks.NewMockAuthService(ctrl)
				s.EXPECT().Decode("valid-token").Return(claim, nil)
				return s
			},
			expStatusCode: http.StatusOK,
		},
		{
			desc:  "should return unauthorized error",
			token: "invalid-token",
			mockSvc: func() *mocks.MockAuthService {
				s := mocks.NewMockAuthService(ctrl)
				s.EXPECT().Decode("invalid-token").Return(nil, errors.New("unauthorized"))
				return s
			},
			expStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			s := NewServer(":8080", nil, tc.mockSvc())

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/authorization/validate?token=%s", tc.token), nil)

			router := mux.NewRouter()
			router.Methods(http.MethodGet).Path("/api/v1/authorization/validate").HandlerFunc(s.tokenValidation)
			router.ServeHTTP(w, r)
			assert.Equal(t, tc.expStatusCode, w.Code)
		})
	}
}
