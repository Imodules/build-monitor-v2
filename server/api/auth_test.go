package api_test

import (
	"testing"

	"build-monitor-v2/server/api"
	"build-monitor-v2/server/db"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2/bson"
)

func TestGenerateAndValidateToken(t *testing.T) {

	Convey("Given a valid user and secret", t, func() {
		secret := "Shh! Don't Tell"
		user := db.User{
			DbObject: db.DbObject{Id: bson.NewObjectId()},
			Username: "the cool 2",
			Email:    "pstuart@fwe.com",
		}

		signedToken, err := api.GenerateToken(&user, secret)

		Convey("It should return a valid signed token", func() {
			So(err, ShouldBeNil)
			So(signedToken, ShouldNotBeEmpty)
		})

		Convey("It should be able to get the token from the context and be valid", func() {
			claims, err := getClaims(signedToken, secret)

			So(err, ShouldBeNil)
			So(claims, ShouldNotBeNil)

			So(claims.UserId, ShouldEqual, user.Id.Hex())
			So(claims.Username, ShouldEqual, user.Username)
			So(claims.Email, ShouldEqual, user.Email)
			So(claims.Valid(), ShouldBeNil)

		})
	})

}
