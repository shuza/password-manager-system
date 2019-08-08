package api

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"password-service/error_tracer"
	"password-service/service"
	"testing"
)

//	Only test unauthorized access as authorized access will be covered by other test suits
func TestAuthorization(t *testing.T) {
	error_tracer.Client = &error_tracer.MockLog{}

	mockAuth := service.AuthMock{}
	mockAuth.On("GetUserId", "invalid-token").Return(0, errors.New("invalid token"))
	service.AuthService = &mockAuth

	Convey("Request password list with invalid token", t, func() {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/password", nil)
		req.Header.Add("Authorization", "invalid-token")
		resp := httptest.NewRecorder()
		NewGinEngine().ServeHTTP(resp, req)

		Convey("Should restrict access", func() {
			So(resp.Code, ShouldEqual, http.StatusUnauthorized)
		})
	})
}
