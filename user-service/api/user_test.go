package api

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-service/db"
	"user-service/error_tracer"
	"user-service/model"
)

func Test_createUserInvalidRequest(t *testing.T) {
	Convey("Request create user api with invalid request body", t, func() {
		db.Client = &db.MockDb{}
		error_tracer.Client = &error_tracer.MockLog{}

		req := httptest.NewRequest(http.MethodPost, "/api/v1/user", bytes.NewBuffer([]byte("invalid request body")))
		resp := httptest.NewRecorder()
		NewGinEngine().ServeHTTP(resp, req)

		Convey("Should return bad request response", func() {
			So(resp.Code, ShouldEqual, http.StatusBadRequest)
		})
	})

	Convey("Request create user api with empty required fields", t, func() {
		db.Client = &db.MockDb{}
		error_tracer.Client = &error_tracer.MockLog{}

		requestBody, _ := json.Marshal(map[string]string{
			"email":         "",
			"password":      "",
			"fullname":      "",
			"business_name": "",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/user", bytes.NewBuffer(requestBody))
		resp := httptest.NewRecorder()
		NewGinEngine().ServeHTTP(resp, req)

		Convey("Should return bad request response", func() {
			So(resp.Code, ShouldEqual, http.StatusBadRequest)
		})
	})

}

func Test_createUserInternalServerError(t *testing.T) {
	Convey("Request create api for internal server error", t, func() {
		error_tracer.Client = &error_tracer.MockLog{}
		user := model.User{
			Email:        "user-1@gmail.com",
			Password:     "123456",
			Fullname:     "Mr user-1",
			BusinessName: "Business-1",
		}
		requestBody, _ := json.Marshal(user)

		h := md5.New()
		h.Write([]byte(user.Password))
		user.Password = hex.EncodeToString(h.Sum(nil))

		mockDb := db.MockDb{}
		mockDb.On("Save", &user).Return(errors.New("Invalid data"))
		db.Client = &mockDb

		req := httptest.NewRequest(http.MethodPost, "/api/v1/user", bytes.NewBuffer(requestBody))
		resp := httptest.NewRecorder()
		NewGinEngine().ServeHTTP(resp, req)

		Convey("Should return internal server error", func() {
			So(resp.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})
}

func Test_createUserSuccess(t *testing.T) {
	Convey("Request create api for success", t, func() {
		error_tracer.Client = &error_tracer.MockLog{}
		user := model.User{
			Email:        "user-1@gmail.com",
			Password:     "123456",
			Fullname:     "Mr user-1",
			BusinessName: "Business-1",
		}
		requestBody, _ := json.Marshal(user)

		h := md5.New()
		h.Write([]byte(user.Password))
		user.Password = hex.EncodeToString(h.Sum(nil))

		mockDb := db.MockDb{}
		mockDb.On("Save", &user).Return(nil)
		db.Client = &mockDb

		req := httptest.NewRequest(http.MethodPost, "/api/v1/user", bytes.NewBuffer(requestBody))
		resp := httptest.NewRecorder()
		NewGinEngine().ServeHTTP(resp, req)

		Convey("Should return internal server error", func() {
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
	})
}

func Test_loginInvalidRequest(t *testing.T) {
	Convey("Request login api with invalid request body", t, func() {
		db.Client = &db.MockDb{}
		error_tracer.Client = &error_tracer.MockLog{}

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer([]byte("invalid request body")))
		resp := httptest.NewRecorder()
		NewGinEngine().ServeHTTP(resp, req)

		Convey("Should return bad request response", func() {
			So(resp.Code, ShouldEqual, http.StatusBadRequest)
		})
	})
}

func Test_loginFail(t *testing.T) {
	Convey("Request login api with non exist user", t, func() {
		error_tracer.Client = &error_tracer.MockLog{}
		user := model.User{
			Email:    "user-1@gmail.com",
			Password: "123456",
		}
		requestBody, _ := json.Marshal(user)

		mockDb := db.MockDb{}
		mockDb.On("GetByEmail", user.Email).Return(user, errors.New("User should not exist"))
		db.Client = &mockDb

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(requestBody))
		resp := httptest.NewRecorder()
		NewGinEngine().ServeHTTP(resp, req)

		Convey("Should not find any user", func() {
			So(resp.Code, ShouldEqual, http.StatusNotFound)
		})
	})

	Convey("Request login api with non exist user", t, func() {
		error_tracer.Client = &error_tracer.MockLog{}
		user := model.User{
			Email:    "user-1@gmail.com",
			Password: "123456",
		}
		requestBody, _ := json.Marshal(user)

		mockDb := db.MockDb{}
		mockDb.On("GetByEmail", user.Email).Return(user, nil)
		db.Client = &mockDb

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(requestBody))
		resp := httptest.NewRecorder()
		NewGinEngine().ServeHTTP(resp, req)

		Convey("Credential Should not match", func() {
			So(resp.Code, ShouldEqual, http.StatusNotAcceptable)
		})
	})
}

func Test_loginSuccess(t *testing.T) {
	Convey("Request login api with valid credential", t, func() {
		error_tracer.Client = &error_tracer.MockLog{}
		user := model.User{
			Email:    "user-1@gmail.com",
			Password: "123456",
		}
		requestBody, _ := json.Marshal(user)

		h := md5.New()
		h.Write([]byte(user.Password))
		user.Password = hex.EncodeToString(h.Sum(nil))

		mockDb := db.MockDb{}
		mockDb.On("GetByEmail", user.Email).Return(user, nil)
		db.Client = &mockDb

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(requestBody))
		resp := httptest.NewRecorder()
		NewGinEngine().ServeHTTP(resp, req)

		Convey("Should login successfully", func() {
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
	})
}
