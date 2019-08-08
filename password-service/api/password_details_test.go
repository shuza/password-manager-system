package api

import (
	"bytes"
	"encoding/json"
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"password-service/db"
	"password-service/error_tracer"
	"password-service/model"
	"password-service/service"
	"testing"
)

func TestPasswordDetails(t *testing.T) {
	error_tracer.Client = &error_tracer.MockLog{}

	mockAuth := service.AuthMock{}
	mockAuth.On("GetUserId", "token-1").Return(1, nil)
	service.AuthService = &mockAuth

	mockDb := db.MockDb{}
	password := model.Password{
		UserId:      1,
		Email:       "user-1@gmail.com",
		Password:    "123456",
		AccountName: "account-1",
		Username:    "username-1",
	}
	mockDb.On("GetById", 0, 1).Return(model.Password{}, errors.New("invalid password id"))
	mockDb.On("GetById", 1, 1).Return(password, nil)
	db.Client = &mockDb

	Convey("Request for password details with non exist id", t, func() {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/password/0", nil)
		req.Header.Add("Authorization", "token-1")
		resp := httptest.NewRecorder()

		NewGinEngine().ServeHTTP(resp, req)
		Convey("Should return not found response", func() {
			So(resp.Code, ShouldEqual, http.StatusNotFound)
		})
	})

	Convey("Request for password details with valid id", t, func() {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/password/1", nil)
		req.Header.Add("Authorization", "token-1")
		resp := httptest.NewRecorder()

		NewGinEngine().ServeHTTP(resp, req)
		Convey("Should return password details", func() {
			So(resp.Code, ShouldEqual, http.StatusOK)

			data, _ := ioutil.ReadAll(resp.Body)
			var response struct {
				Status string         `json:"status"`
				Data   model.Password `json:"data"`
			}
			json.Unmarshal(data, &response)

			So(response.Data.Password, ShouldEqual, password.Password)
		})
	})
}

func TestDeletePassword(t *testing.T) {
	error_tracer.Client = &error_tracer.MockLog{}

	mockAuth := service.AuthMock{}
	mockAuth.On("GetUserId", "token-1").Return(1, nil)
	mockAuth.On("GetUserId", "token-2").Return(2, nil)
	service.AuthService = &mockAuth

	mockDb := db.MockDb{}
	password := model.Password{
		UserId:      1,
		Email:       "user-1@gmail.com",
		Password:    "123456",
		AccountName: "account-1",
		Username:    "username-1",
	}
	mockDb.On("GetById", 0, 1).Return(model.Password{}, errors.New("invalid password id"))
	mockDb.On("GetById", 1, 1).Return(password, nil)
	mockDb.On("Delete", &password).Return(errors.New("invalid delete operation"))

	password2 := model.Password{
		UserId:      2,
		Email:       "user-1@gmail.com",
		Password:    "123456",
		AccountName: "account-1",
		Username:    "username-1",
	}
	mockDb.On("GetById", 1, 2).Return(password2, nil)
	mockDb.On("Delete", &password2).Return(nil)
	db.Client = &mockDb

	Convey("Request to delete non exist password id", t, func() {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/password/0", nil)
		req.Header.Add("Authorization", "token-1")
		resp := httptest.NewRecorder()

		NewGinEngine().ServeHTTP(resp, req)
		Convey("Should return not found response", func() {
			So(resp.Code, ShouldEqual, http.StatusNotFound)
		})
	})

	Convey("Request to delete password id but database fail", t, func() {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/password/1", nil)
		req.Header.Add("Authorization", "token-1")
		resp := httptest.NewRecorder()

		NewGinEngine().ServeHTTP(resp, req)
		Convey("Should return not found response", func() {
			So(resp.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})

	Convey("Request to delete valid password", t, func() {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/password/1", nil)
		req.Header.Add("Authorization", "token-2")
		resp := httptest.NewRecorder()

		NewGinEngine().ServeHTTP(resp, req)
		Convey("Should return not found response", func() {
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
	})
}

func TestUpdatePasswordInvalidRequest(t *testing.T) {
	error_tracer.Client = &error_tracer.MockLog{}

	mockAuth := service.AuthMock{}
	mockAuth.On("GetUserId", "token-1").Return(1, nil)
	service.AuthService = &mockAuth

	mockDb := db.MockDb{}
	password := model.Password{
		UserId:      1,
		Email:       "user-1@gmail.com",
		Password:    "123456",
		AccountName: "account-1",
		Username:    "username-1",
	}
	mockDb.On("GetById", 0, 1).Return(model.Password{}, errors.New("invalid password id"))
	mockDb.On("GetById", 1, 1).Return(password, nil)
	db.Client = &mockDb

	requestBody, _ := json.Marshal(password)

	Convey("Request to update password with invalid request body", t, func() {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/password/0", bytes.NewBuffer([]byte("invalid request body")))
		req.Header.Add("Authorization", "token-1")
		resp := httptest.NewRecorder()
		NewGinEngine().ServeHTTP(resp, req)

		Convey("Should return bad request response", func() {
			So(resp.Code, ShouldEqual, http.StatusBadRequest)
		})
	})

	Convey("Request to update password with non exist id", t, func() {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/password/0", bytes.NewBuffer(requestBody))
		req.Header.Add("Authorization", "token-1")
		req.Header.Add("user_id", "1")
		resp := httptest.NewRecorder()
		NewGinEngine().ServeHTTP(resp, req)

		Convey("Should return bad request response", func() {
			So(resp.Code, ShouldEqual, http.StatusNotFound)
		})
	})
}

func TestUpdatePasswordDatabaseInteraction(t *testing.T) {
	error_tracer.Client = &error_tracer.MockLog{}

	mockAuth := service.AuthMock{}
	mockAuth.On("GetUserId", "token-1").Return(1, nil)
	mockAuth.On("GetUserId", "token-2").Return(2, nil)
	service.AuthService = &mockAuth

	mockDb := db.MockDb{}
	password := model.Password{
		UserId:      1,
		Email:       "user-1@gmail.com",
		Password:    "123456",
		AccountName: "account-1",
		Username:    "username-1",
	}
	mockDb.On("GetById", 1, 1).Return(password, nil)
	mockDb.On("Save", &password).Return(errors.New("invalid dave request"))

	password2 := model.Password{
		UserId:      2,
		Email:       "user-2",
		Password:    "123456",
		AccountName: "account-2",
		Username:    "username-2",
	}
	mockDb.On("GetById", 2, 2).Return(password2, nil)
	mockDb.On("Save", &password2).Return(nil)
	db.Client = &mockDb

	Convey("Request to update password but database fail", t, func() {
		requestBody, _ := json.Marshal(password)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/password/1", bytes.NewBuffer(requestBody))
		req.Header.Add("Authorization", "token-1")
		resp := httptest.NewRecorder()
		NewGinEngine().ServeHTTP(resp, req)

		Convey("Should return internal server error response", func() {
			So(resp.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})

	Convey("Request to update password successful", t, func() {
		requestBody, _ := json.Marshal(password2)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/password/2", bytes.NewBuffer(requestBody))
		req.Header.Add("Authorization", "token-2")
		resp := httptest.NewRecorder()
		NewGinEngine().ServeHTTP(resp, req)

		Convey("Should return successful response", func() {
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
	})
}
