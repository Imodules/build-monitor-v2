package api_test

import (
	"testing"

	"build-monitor-v2/server/api"
	"build-monitor-v2/server/cfg"

	"errors"

	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2"
)

func TestCreate(t *testing.T) {

	Convey("Given a logger, config", t, func() {
		conf := cfg.Config{}
		log := logrus.WithField("test", "TestCreate")
		session := mgo.Session{}

		Convey("It should return a new s object", func() {
			server := api.Create(log, &conf, &session)

			So(server, ShouldNotBeNil)
			So(server.Config, ShouldEqual, &conf)
			So(server.Log, ShouldEqual, log)
			So(server.DbSession, ShouldEqual, &session)
			So(server.Server, ShouldNotBeNil)
		})

	})

}

func TestStart(t *testing.T) {
	Convey("Given a server", t, func() {

		serverMock := new(IServerMock)

		server := api.Server{
			Server: serverMock,
			Config: &cfg.Config{Port: 10101},
		}

		Convey("When Start is called", func() {
			serverMock.On("Start", ":10101").Return(nil)

			server.Start()

			Convey("It should start the server with the proper port", func() {
				serverMock.AssertExpectations(t)
			})

		})

		Convey("When server.Start has an error", func() {
			expectedError := errors.New("This is the error!")
			serverMock.On("Start", ":10101").Return(expectedError)

			err := server.Start()

			Convey("It should return that error", func() {
				So(err, ShouldEqual, expectedError)
			})
		})

	})
}

func TestServer_Shutdown(t *testing.T) {
	Convey("Given a server", t, func() {
		serverMock := new(IServerMock)
		log := logrus.WithField("test", "TestServer_Shutdown")

		server := api.Server{
			Server: serverMock,
			Config: &cfg.Config{Port: 10101},
			Log:    log,
		}

		Convey("When shutdown is called successfully", func() {
			serverMock.On("Shutdown", mock.AnythingOfType("*context.timerCtx")).Return(nil)

			server.Shutdown()

			Convey("It should not panic", func() {
				serverMock.AssertExpectations(t)
			})
		})

		Convey("When shutdown panics", func() {
			serverMock.On("Shutdown", mock.AnythingOfType("*context.timerCtx")).Return(errors.New("panicking"))

			So(func() { server.Shutdown() }, ShouldPanic)
		})
	})
}
