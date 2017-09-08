package cfg_test

import (
	"testing"

	"build-monitor-v2/server/cfg"

	"errors"

	. "github.com/smartystreets/goconvey/convey"
)

func runConfigAsserts(c *cfg.Config) {
	So(c.Db, ShouldEqual, "mongodb://localhost:27017/build-monitor-v2")
	So(c.Port, ShouldEqual, 3030)
	So(c.ClientPath, ShouldEqual, "../client/dist")
	So(c.AllowedOrigin, ShouldEqual, "*")
	So(c.PasswordSalt, ShouldEqual, "you-really-need-to-change-this")
	So(c.JwtSecret, ShouldEqual, "you-really-need-to-change-this-one-also")
}

func TestLoad(t *testing.T) {
	Convey("Given a successful getOverrides func", t, func() {
		var overridesCfg *cfg.Config

		var getOverrides = func(s interface{}) error {
			overridesCfg = s.(*cfg.Config)

			return nil
		}

		Convey("When Load is called", func() {
			result, err := cfg.Load(getOverrides)

			Convey("It should call getOverrides with a default config", func() {

				runConfigAsserts(overridesCfg)

				Convey("And return a fully populated config with defaults", func() {

					So(result, ShouldNotBeNil)

					runConfigAsserts(&result)

					Convey("And a nil error", func() {
						So(err, ShouldBeNil)
					})
				})
			})

		})

	})

	Convey("Given an unsuccessful getOverrides func", t, func() {
		var getOverrides = func(s interface{}) error {
			return errors.New("This was not a success!")
		}

		Convey("When Load is called", func() {
			result, err := cfg.Load(getOverrides)

			Convey("It should return an error with the default config", func() {

				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "This was not a success!")

				So(result, ShouldNotBeNil)

				runConfigAsserts(&result)
			})

		})

	})
}
