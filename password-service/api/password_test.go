package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

func TestAddPasswordRequestBody(t *testing.T) {
	mockAuth := service.AuthMock{}
	mockAuth.On("GetUserId", "").Return(1, nil)
	service.AuthService = &mockAuth

	Convey("Request add password api with wrong request body", t, func() {
		error_tracer.Client = &error_tracer.MockLog{}
		req := httptest.NewRequest(http.MethodPost, "/api/v1/password", bytes.NewBuffer([]byte("invalid request body")))
		resp := httptest.NewRecorder()
		NewGinEngine().ServeHTTP(resp, req)

		Convey("Should return bad request response", func() {
			So(resp.Code, ShouldEqual, http.StatusBadRequest)
		})
	})

	Convey("Request add password api with invalid password", t, func() {
		error_tracer.Client = &error_tracer.MockLog{}
		requestBody, _ := json.Marshal(model.Password{})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/password", bytes.NewBuffer(requestBody))
		resp := httptest.NewRecorder()
		NewGinEngine().ServeHTTP(resp, req)

		Convey("Should return bad request response", func() {
			So(resp.Code, ShouldEqual, http.StatusBadRequest)
		})
	})
}

func TestAddPasswordWithoutToken(t *testing.T) {
	error_tracer.Client = &error_tracer.MockLog{}
	mockAuth := service.AuthMock{}
	mockAuth.On("GetUserId", "").Return(0, nil)
	service.AuthService = &mockAuth

	Convey("Request add password without token in header", t, func() {
		requestBody, _ := json.Marshal(model.Password{
			Password:    "123456",
			Username:    "aaa",
			AccountName: "account",
			Email:       "aaa@aaa.com",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/password", bytes.NewBuffer(requestBody))
		resp := httptest.NewRecorder()
		NewGinEngine().ServeHTTP(resp, req)

		Convey("Should return internal server error response", func() {
			So(resp.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})
}

func TestAddPasswordDatabaseInteraction(t *testing.T) {
	error_tracer.Client = &error_tracer.MockLog{}

	mockAuth := service.AuthMock{}
	mockAuth.On("GetUserId", "valid-token").Return(1, nil)
	service.AuthService = &mockAuth

	data := model.Password{
		UserId:      1,
		Password:    "123456",
		Username:    "aaa",
		AccountName: "account",
		Email:       "aaa@aaa.com",
	}
	requestBody, _ := json.Marshal(data)

	Convey("Request add password for database fail", t, func() {
		mockDb := db.MockDb{}
		mockDb.On("Save", &data).Return(errors.New("Database fail test for add password"))
		db.Client = &mockDb

		req := httptest.NewRequest(http.MethodPost, "/api/v1/password", bytes.NewBuffer(requestBody))
		req.Header.Add("Authorization", "valid-token")
		resp := httptest.NewRecorder()
		NewGinEngine().ServeHTTP(resp, req)

		Convey("Should return internal server error response", func() {
			So(resp.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})

	Convey("Request add password for database fail", t, func() {
		mockDb := db.MockDb{}
		mockDb.On("Save", &data).Return(nil)
		db.Client = &mockDb

		req := httptest.NewRequest(http.MethodPost, "/api/v1/password", bytes.NewBuffer(requestBody))
		req.Header.Add("Authorization", "valid-token")
		resp := httptest.NewRecorder()
		NewGinEngine().ServeHTTP(resp, req)

		Convey("Should return internal server error response", func() {
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
	})
}

func TestPasswordList(t *testing.T) {
	error_tracer.Client = &error_tracer.MockLog{}

	mockAuth := service.AuthMock{}
	mockAuth.On("GetUserId", "valid-token").Return(1, nil)
	mockAuth.On("GetUserId", "valid-token-2").Return(2, nil)
	service.AuthService = &mockAuth

	mockDb := db.MockDb{}
	mockDb.On("GetByUserId", uint(2)).Return([]model.Password{}, errors.New("invalid user_id"))

	passwords := getDummyPasswordList(4, 1)
	mockDb.On("GetByUserId", uint(1)).Return(passwords, nil)
	db.Client = &mockDb

	Convey("Request for password list without user_id", t, func() {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/password", nil)
		req.Header.Add("Authorization", "valid-token-2")
		resp := httptest.NewRecorder()
		NewGinEngine().ServeHTTP(resp, req)

		Convey("Should return internal server error response", func() {
			So(resp.Code, ShouldEqual, http.StatusInternalServerError)
		})

	})

	Convey("Request for password list without user_id", t, func() {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/password", nil)
		req.Header.Add("Authorization", "valid-token")
		resp := httptest.NewRecorder()
		NewGinEngine().ServeHTTP(resp, req)

		Convey("Should return 4 password", func() {
			So(resp.Code, ShouldEqual, http.StatusOK)
			data, _ := ioutil.ReadAll(resp.Body)
			var response struct {
				Status string           `json:"status"`
				Data   []model.Password `json:"data"`
			}
			json.Unmarshal(data, &response)
			So(len(response.Data), ShouldEqual, 4)
		})

	})
}

func getDummyPasswordList(n int, userId int) []model.Password {
	passwords := make([]model.Password, 0)
	for i := 1; i <= n; i++ {
		password := model.Password{
			UserId:      userId,
			Password:    fmt.Sprintf("passoword-%d", i),
			AccountName: fmt.Sprintf("account-%d", i),
			Username:    fmt.Sprintf("username-%d", i),
			Email:       fmt.Sprintf("email-%d@email.com", i),
		}

		passwords = append(passwords, password)
	}

	return passwords
}
