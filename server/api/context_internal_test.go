package api

import (
	"testing"

	"net/http"
	"net/http/httptest"
	"strings"

	"errors"

	"build-monitor-v2/server/cfg"

	"build-monitor-v2/server/db"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2"
)

func TestContextGetters(t *testing.T) {

	Convey("Given an echo.Context", t, func() {
		e := echo.New()
		req, _ := http.NewRequest(echo.GET, "/some/path", strings.NewReader(""))
		rec := httptest.NewRecorder()

		Convey("When that context does not have a logger", func() {
			c := e.NewContext(req, rec)
			foundLogger := getLogger(c)

			Convey("getLogger should return the logger", func() {
				So(foundLogger, ShouldNotBeNil)
			})
		})

		Convey("When that context has a logger", func() {
			c := e.NewContext(req, rec)

			logger := logrus.WithField("test", "TestGetLogger1")
			c.Set(loggerKey, logger)

			foundLogger := getLogger(c)

			Convey("getLogger should return the logger", func() {
				So(foundLogger, ShouldNotBeNil)
				So(foundLogger, ShouldEqual, logger)
			})
		})

		Convey("When that context does not have a token", func() {
			c := e.NewContext(req, rec)
			foundToken := getToken(c)

			Convey("getToken should return nil", func() {
				So(foundToken, ShouldBeNil)
			})

		})

		Convey("When that context has a token", func() {
			c := e.NewContext(req, rec)
			token := &jwt.Token{}
			c.Set(tokenKey, token)

			foundToken := getToken(c)

			Convey("getToken should return that token", func() {
				So(foundToken, ShouldEqual, token)
			})

		})

		Convey("When that context does not have a dbSession", func() {
			c := e.NewContext(req, rec)
			foundDbSession := getAppDb(c)

			Convey("getAppDb should return nil", func() {
				So(foundDbSession, ShouldBeNil)
			})

		})

		Convey("When that context has a dbSession", func() {
			c := e.NewContext(req, rec)

			Db := db.AppDb{PasswordSalt: "something cool and diff"}
			c.Set(dbKey, &Db)

			foundDbSession := getAppDb(c)

			Convey("getAppDb should return the Db", func() {
				So(foundDbSession, ShouldPointTo, &Db)
			})

		})
	})

}

func TestGetSetupRequestHandler(t *testing.T) {
	session, _ := mgo.Dial("mongodb://localhost/build-monitor-v2-test")
	defer session.Close()

	Convey("Given an API object", t, func() {
		server := Server{
			Config:    &cfg.Config{Port: 7630, PasswordSalt: "dumm fake one"},
			Log:       logrus.WithField("test", "TestSetupRequest"),
			DbSession: session,
		}

		Convey("When the request is successful", func() {
			e := echo.New()
			req, _ := http.NewRequest(echo.GET, "/users", strings.NewReader(""))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			var mockHandlerContext echo.Context
			var mockHandler = func(ctx echo.Context) error {
				mockHandlerContext = ctx
				return nil
			}

			handlerError := getSetupRequestHandler(&server)(mockHandler)(c)

			So(c.Get(loggerKey), ShouldNotBeNil)
			So(c.Get(dbKey), ShouldNotBeNil)

			Convey("The error should not be set", func() {
				So(mockHandlerContext, ShouldEqual, c)
				So(handlerError, ShouldBeNil)
				So(rec.Code, ShouldEqual, http.StatusOK)
			})
		})

		Convey("When the request has an error", func() {
			expectedError := errors.New("Something went wrong")

			e := echo.New()
			req, _ := http.NewRequest(echo.GET, "/users", strings.NewReader(""))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			var mockHandlerContext echo.Context
			var mockHandler = func(ctx echo.Context) error {
				mockHandlerContext = ctx
				return expectedError
			}

			handlerError := getSetupRequestHandler(&server)(mockHandler)(c)

			So(c.Get(loggerKey), ShouldNotBeNil)
			So(c.Get(dbKey), ShouldNotBeNil)

			Convey("The error should not be set", func() {
				So(mockHandlerContext, ShouldEqual, c)
				So(handlerError, ShouldEqual, expectedError)
				So(rec.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}
