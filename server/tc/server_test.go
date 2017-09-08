package tc_test

import (
	"testing"

	"build-monitor-v2/server/cfg"

	"build-monitor-v2/server/tc"

	"time"

	"errors"
	"github.com/kapitanov/go-teamcity"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewServer(t *testing.T) {

	Convey("Given a logger, config", t, func() {
		Convey("When the intervals are valid", func() {
			conf := cfg.Config{
				TcProjectPollInterval:      "2h45m",
				TcBuildPollInterval:        "1500ms",
				TcRunningBuildPollInterval: "2m",
			}
			log := logrus.WithField("test", "TestNewServer")

			Convey("It should return a new s object", func() {
				server := tc.NewServer(log, &conf)

				So(server, ShouldNotBeNil)
				So(server.Log, ShouldEqual, log)
				So(server.ProjectPollInterval, ShouldEqual, (time.Hour*2)+(time.Minute*45))
				So(server.BuildPollInterval, ShouldEqual, time.Millisecond*1500)
				So(server.RunningBuildPollInterval, ShouldEqual, time.Minute*2)
			})
		})

		Convey("When the intervals are not valid", func() {
			conf := cfg.Config{
				TcProjectPollInterval:      "a",
				TcBuildPollInterval:        "b",
				TcRunningBuildPollInterval: "c",
			}
			log := logrus.WithField("test", "TestNewServer")

			Convey("It should return a new s object", func() {
				So(func() { tc.NewServer(log, &conf) }, ShouldPanic)
			})
		})
	})

}

func TestServer_Start_Shutdown(t *testing.T) {
	Convey("Given a logger and tcServer", t, func() {
		log := logrus.WithField("test", "TestServer_Start_Shutdown")
		serverMock := new(ITcClientMock)

		c := tc.Server{
			Tc:                  serverMock,
			Log:                 log,
			ProjectPollInterval: time.Millisecond * 500,
		}

		projects := []teamcity.Project{}

		Convey("When we successfully start the monitor", func() {
			Convey("It should call GetProjects at startup", func() {

				serverMock.On("GetProjects").Times(2).Return(projects, nil)

				err := c.Start()
				So(err, ShouldBeNil)
				So(len(serverMock.Calls), ShouldEqual, 1)

				Convey("And call again after the timeout", func() {
					<-time.After(time.Millisecond * 750)

					So(len(serverMock.Calls), ShouldEqual, 2)

					c.Shutdown()

					serverMock.AssertExpectations(t)
				})

			})
		})

		Convey("When we fail to start the monitor", func() {
			Convey("It should call GetProjects at startup", func() {

				expectedError := errors.New("there was something wrong")
				serverMock.On("GetProjects").Times(1).Return(projects, expectedError)

				err := c.Start()
				So(err, ShouldEqual, expectedError)

				serverMock.AssertExpectations(t)

			})
		})
	})
}
