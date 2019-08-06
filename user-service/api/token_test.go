package api

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-service/error_tracer"
	"user-service/model"
	"user-service/service"
)

func Test_tokenVerify(t *testing.T) {
	Convey("Request token validation api with invalid token", t, func() {
		error_tracer.Client = &error_tracer.MockLog{}
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/token?token=no-token", nil)
		resp := httptest.NewRecorder()
		NewGinEngine().ServeHTTP(resp, req)

		Convey("Should restrict access", func() {
			So(resp.Code, ShouldEqual, http.StatusUnauthorized)
		})
	})

	Convey("Request token validation api with invalid token", t, func() {
		error_tracer.Client = &error_tracer.MockLog{}

		tokenService := service.TokenService{}
		token, err := tokenService.Encode(model.User{
			Email:        "user-1@gmail.com",
			Password:     "123456",
			Fullname:     "Mr. User-1",
			BusinessName: "Business-1",
		})
		if err != nil {
			panic(err)
		}

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/auth/token?token=%s", token), nil)
		resp := httptest.NewRecorder()
		NewGinEngine().ServeHTTP(resp, req)

		Convey("Should grant access", func() {
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
	})
}
