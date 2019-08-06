package api

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_index(t *testing.T) {
	Convey("Test index api always ok", t, func() {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		resp := httptest.NewRecorder()
		NewGinEngine().ServeHTTP(resp, req)

		Convey("Should always success", func() {
			So(resp.Code, ShouldEqual, http.StatusOK)
		})
	})
}
