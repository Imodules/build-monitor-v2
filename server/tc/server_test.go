package tc_test

import (
	"testing"

	"build-monitor-v2/server/cfg"

	"build-monitor-v2/server/tc"

	"time"

	"errors"

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

			dbMock := new(IDbMock)

			Convey("It should return a new s object", func() {
				server := tc.NewServer(log, &conf, dbMock)

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
			dbMock := new(IDbMock)

			Convey("It should return a new s object", func() {
				So(func() { tc.NewServer(log, &conf, dbMock) }, ShouldPanic)
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

		Convey("When we successfully start the monitor", func() {
			oldRefreshProjects := tc.RefreshProjects
			refreshProjectsCallCount := 0
			tc.RefreshProjects = func(tcs *tc.Server) error {
				refreshProjectsCallCount++
				return nil
			}
			defer func() { tc.RefreshProjects = oldRefreshProjects }()

			oldRefreshBuildTypes := tc.RefreshBuildTypes
			refreshBuildTypesCallCount := 0
			tc.RefreshBuildTypes = func(tcs *tc.Server) error {
				refreshBuildTypesCallCount++
				return nil
			}
			defer func() { tc.RefreshBuildTypes = oldRefreshBuildTypes }()

			Convey("It should call GetProjects at startup", func() {

				err := c.Start()
				So(err, ShouldBeNil)
				So(refreshProjectsCallCount, ShouldEqual, 1)
				So(refreshBuildTypesCallCount, ShouldEqual, 1)

				Convey("And call again after the timeout", func() {
					<-time.After(time.Millisecond * 750)

					So(refreshProjectsCallCount, ShouldEqual, 2)
					So(refreshBuildTypesCallCount, ShouldEqual, 2)

					c.Shutdown()
				})

			})
		})

		Convey("When we fail to start the monitor because RefreshProjects fails", func() {
			expectedError := errors.New("there was something wrong")

			oldRefreshProjects := tc.RefreshProjects
			refreshProjectsCallCount := 0
			tc.RefreshProjects = func(tcs *tc.Server) error {
				refreshProjectsCallCount++
				return expectedError
			}
			defer func() { tc.RefreshProjects = oldRefreshProjects }()

			oldRefreshBuildTypes := tc.RefreshBuildTypes
			refreshBuildTypesCallCount := 0
			tc.RefreshBuildTypes = func(tcs *tc.Server) error {
				refreshBuildTypesCallCount++
				return nil
			}
			defer func() { tc.RefreshBuildTypes = oldRefreshBuildTypes }()

			Convey("It should call RefreshProjects at startup", func() {

				err := c.Start()
				So(err, ShouldEqual, expectedError)
				So(refreshProjectsCallCount, ShouldEqual, 1)
				So(refreshBuildTypesCallCount, ShouldEqual, 0)
			})
		})

		Convey("When we fail to start the monitor because RefreshBuildTypes fails", func() {
			expectedError := errors.New("there was something wrong")

			oldRefreshProjects := tc.RefreshProjects
			refreshProjectsCallCount := 0
			tc.RefreshProjects = func(tcs *tc.Server) error {
				refreshProjectsCallCount++
				return nil
			}
			defer func() { tc.RefreshProjects = oldRefreshProjects }()

			oldRefreshBuildTypes := tc.RefreshBuildTypes
			refreshBuildTypesCallCount := 0
			tc.RefreshBuildTypes = func(tcs *tc.Server) error {
				refreshBuildTypesCallCount++
				return expectedError
			}
			defer func() { tc.RefreshBuildTypes = oldRefreshBuildTypes }()

			Convey("It should call RefreshProjects at startup", func() {

				err := c.Start()
				So(err, ShouldEqual, expectedError)
				So(refreshProjectsCallCount, ShouldEqual, 1)
				So(refreshBuildTypesCallCount, ShouldEqual, 1)
			})
		})
	})
}
