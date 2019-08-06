package service

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"user-service/model"
)

func TestTokenService(t *testing.T) {
	tokenService := TokenService{}
	user := model.User{
		Email:        "user-1@gmail.com",
		Fullname:     "Mr. User-1",
		Password:     "123456",
		BusinessName: "Business-1",
	}

	Convey("Encode a user then decode it and check with decoded data", t, func() {
		tokenStr, err := tokenService.Encode(user)
		Convey("Should encode successfully", func() {
			So(err, ShouldBeNil)
			So(tokenStr, ShouldNotBeEmpty)
		})

		decodedUser, err := tokenService.Decode(tokenStr)
		Convey("Should decode exact user data", func() {
			So(err, ShouldBeNil)
			So(decodedUser.User.Email, ShouldEqual, user.Email)
			So(decodedUser.User.Fullname, ShouldEqual, user.Fullname)
			So(decodedUser.User.Password, ShouldEqual, user.Password)
			So(decodedUser.User.BusinessName, ShouldEqual, user.BusinessName)
		})
	})
}

func TestTokenServiceDecode(t *testing.T) {
	tokenService := TokenService{}
	Convey("Test decoder with invalid token", t, func() {
		Convey("Should arise error", func() {
			_, err := tokenService.Decode("abcdefghij1234567890")
			So(err, ShouldNotBeNil)
		})
	})
}
