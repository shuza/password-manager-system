package api

import (
	"encoding/json"
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"password-service/db"
	"password-service/error_tracer"
	"password-service/model"
	"testing"
)

func TestPasswordDetails(t *testing.T) {
	error_tracer.Client = &error_tracer.MockLog{}
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
		req.Header.Add("user_id", "1")
		resp := httptest.NewRecorder()

		NewGinEngine().ServeHTTP(resp, req)
		Convey("Should return not found response", func() {
			So(resp.Code, ShouldEqual, http.StatusNotFound)
		})
	})

	Convey("Request for password details with valid id", t, func() {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/password/1", nil)
		req.Header.Add("user_id", "1")
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

	password.UserId = 2
	mockDb.On("GetById", 2, 1).Return(password, nil)
	mockDb.On("Delete", &password).Return(nil)
	db.Client = &mockDb

	Convey("Request to delete non exist password id", t, func() {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/password/0", nil)
		req.Header.Add("user_id", "1")
		resp := httptest.NewRecorder()

		NewGinEngine().ServeHTTP(resp, req)
		Convey("Should return not found response", func() {
			So(resp.Code, ShouldEqual, http.StatusNotFound)
		})
	})

	Convey("Request to delete password id but database fail", t, func() {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/password/1", nil)
		req.Header.Add("user_id", "1")
		resp := httptest.NewRecorder()

		NewGinEngine().ServeHTTP(resp, req)
		Convey("Should return not found response", func() {
			So(resp.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})

	/*Convey("Request to delete valid password", t, func() {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/password/1", nil)
		req.Header.Add("user_id", "2")
		resp := httptest.NewRecorder()

		NewGinEngine().ServeHTTP(resp, req)
		Convey("Should return not found response", func() {
			So(resp.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})*/
}
