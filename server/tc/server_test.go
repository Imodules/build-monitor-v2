package tc_test

import (
	"errors"
	"testing"

	"build-monitor-v2/server/cfg"

	"build-monitor-v2/server/tc"

	"time"

	teamcity "github.com/pstuart2/go-teamcity"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewServer(t *testing.T) {

	Convey("Given a logger, config", t, func() {
		Convey("When the intervals are valid", func() {
			conf := cfg.Config{
				TcPollInterval:             "2h45m",
				TcRunningBuildPollInterval: "200ms",
			}
			log := logrus.WithField("test", "TestNewServer")

			dbMock := new(IDbMock)

			Convey("It should return a new s object", func() {
				server := tc.NewServer(log, &conf, dbMock)

				So(server, ShouldNotBeNil)
				So(server.Log, ShouldEqual, log)
				So(server.TcPollInterval, ShouldEqual, (time.Hour*2)+(time.Minute*45))
				So(server.TcRunningBuildPollInterval, ShouldEqual, time.Millisecond*200)
			})
		})

		Convey("When the intervals are not valid", func() {
			conf := cfg.Config{
				TcPollInterval:             "a",
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
			Tc:                         serverMock,
			Log:                        log,
			TcPollInterval:             time.Millisecond * 500,
			TcRunningBuildPollInterval: time.Millisecond * 500,
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

			oldGetBuildHistory := tc.GetBuildHistory
			refreshBuildHistoryCallCount := 0
			tc.GetBuildHistory = func(tcs *tc.Server) error {
				refreshBuildHistoryCallCount++
				return nil
			}
			defer func() { tc.GetBuildHistory = oldGetBuildHistory }()

			Convey("And there are no running builds", func() {
				oldGetRunningBuilds := tc.GetRunningBuilds
				getRunningBuildsCallCount := 0
				tc.GetRunningBuilds = func(tcs *tc.Server, lastBuilds []teamcity.Build) []teamcity.Build {
					getRunningBuildsCallCount++
					return []teamcity.Build{}
				}
				defer func() { tc.GetRunningBuilds = oldGetRunningBuilds }()

				Convey("It should call RefreshProjects at startup", func() {

					err := c.Start()
					So(err, ShouldBeNil)
					So(refreshProjectsCallCount, ShouldEqual, 1)
					So(refreshBuildTypesCallCount, ShouldEqual, 1)
					So(refreshBuildHistoryCallCount, ShouldEqual, 1)
					So(getRunningBuildsCallCount, ShouldEqual, 0)

					Convey("And call GetRunningBuilds after the timeout", func() {
						<-time.After(time.Millisecond * 750)

						So(refreshProjectsCallCount, ShouldEqual, 1)
						So(refreshBuildTypesCallCount, ShouldEqual, 1)
						So(refreshBuildHistoryCallCount, ShouldEqual, 1)
						So(getRunningBuildsCallCount, ShouldEqual, 1)

						c.Shutdown()
					})

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

		// Convey("When we fail to start the monitor because RefreshBuildTypes fails", func() {
		// 	expectedError := errors.New("there was something wrong")

		// 	oldRefreshProjects := tc.RefreshProjects
		// 	refreshProjectsCallCount := 0
		// 	tc.RefreshProjects = func(tcs *tc.Server) error {
		// 		refreshProjectsCallCount++
		// 		return nil
		// 	}
		// 	defer func() { tc.RefreshProjects = oldRefreshProjects }()

		// 	oldRefreshBuildTypes := tc.RefreshBuildTypes
		// 	refreshBuildTypesCallCount := 0
		// 	tc.RefreshBuildTypes = func(tcs *tc.Server) error {
		// 		refreshBuildTypesCallCount++
		// 		return expectedError
		// 	}
		// 	defer func() { tc.RefreshBuildTypes = oldRefreshBuildTypes }()

		// 	Convey("It should call RefreshProjects at startup", func() {

		// 		err := c.Start()
		// 		So(err, ShouldEqual, expectedError)
		// 		So(refreshProjectsCallCount, ShouldEqual, 1)
		// 		So(refreshBuildTypesCallCount, ShouldEqual, 1)
		// 	})
		// })

		// Convey("When we fail to start the monigor because GetBuildHistory fails", func() {
		// 	expectedError := errors.New("there was something wrong")

		// 	oldRefreshProjects := tc.RefreshProjects
		// 	refreshProjectsCallCount := 0
		// 	tc.RefreshProjects = func(tcs *tc.Server) error {
		// 		refreshProjectsCallCount++
		// 		return nil
		// 	}
		// 	defer func() { tc.RefreshProjects = oldRefreshProjects }()

		// 	oldRefreshBuildTypes := tc.RefreshBuildTypes
		// 	refreshBuildTypesCallCount := 0
		// 	tc.RefreshBuildTypes = func(tcs *tc.Server) error {
		// 		refreshBuildTypesCallCount++
		// 		return nil
		// 	}
		// 	defer func() { tc.RefreshBuildTypes = oldRefreshBuildTypes }()

		// 	oldGetBuildHistory := tc.GetBuildHistory
		// 	refreshBuildHistoryCallCount := 0
		// 	tc.GetBuildHistory = func(tcs *tc.Server) error {
		// 		refreshBuildHistoryCallCount++
		// 		return expectedError
		// 	}
		// 	defer func() { tc.GetBuildHistory = oldGetBuildHistory }()

		// 	Convey("It should call RefreshProjects at startup", func() {

		// 		err := c.Start()
		// 		So(err, ShouldEqual, expectedError)
		// 		So(refreshProjectsCallCount, ShouldEqual, 1)
		// 		So(refreshBuildTypesCallCount, ShouldEqual, 1)
		// 		So(refreshBuildHistoryCallCount, ShouldEqual, 1)
		// 	})
		// })
	})
}
