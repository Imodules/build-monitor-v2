package tc_test

import (
	"testing"

	"build-monitor-v2/server/tc"

	"time"

	"build-monitor-v2/server/db"

	"errors"

	"github.com/pstuart2/go-teamcity"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestServer_GetBuildHistory(t *testing.T) {
	Convey("Given a server", t, func() {
		log := logrus.WithField("test", "TestServer_GetBuildHistory")
		serverMock := new(ITcClientMock)
		dbMock := new(IDbMock)

		c := tc.Server{
			Tc:                  serverMock,
			Db:                  dbMock,
			Log:                 log,
			ProjectPollInterval: time.Millisecond * 500,
		}

		Convey("When DashboardList errors", func() {
			expectedErr := errors.New("i knew this would happen")
			dbMock.On("DashboardList").Return(nil, expectedErr)

			err := tc.GetBuildHistory(&c)

			Convey("It should return the error", func() {
				So(err, ShouldEqual, expectedErr)
			})
		})

		Convey("When there are no dashboards", func() {
			dashboards := []db.Dashboard{}

			dbMock.On("DashboardList").Return(dashboards, nil)

			Convey("It should not do anything", func() {
				tc.GetBuildHistory(&c)

				dbMock.AssertExpectations(t)
				serverMock.AssertExpectations(t)
			})
		})

		Convey("When there are dashboards without build configs", func() {
			dashboards := []db.Dashboard{
				{Id: "cool 1"},
				{Id: "cool 2"},
			}

			dbMock.On("DashboardList").Return(dashboards, nil)

			Convey("It should not do anything", func() {
				tc.GetBuildHistory(&c)

				dbMock.AssertExpectations(t)
				serverMock.AssertExpectations(t)
			})
		})

		Convey("When there are dashboards with build configs", func() {
			bcfg1 := db.BuildConfig{Id: "bcfg1"}
			bcfg2 := db.BuildConfig{Id: "bcfg2"}
			bcfg3 := db.BuildConfig{Id: "bcfg3"}
			bcfg4 := db.BuildConfig{Id: "bcfg4"}

			dashboards := []db.Dashboard{
				{Id: "cool 1", BuildConfigs: []db.BuildConfig{bcfg1, bcfg2, bcfg3}},
				{Id: "cool 2", BuildConfigs: []db.BuildConfig{bcfg1, bcfg2, bcfg4}},
			}

			dbMock.On("DashboardList").Return(dashboards, nil)

			Convey("And there are no builds", func() {
				builds := []teamcity.Build{}

				serverMock.On("GetBuildsForBuildType", "bcfg1", 1000).Times(1).Return(builds, nil)
				serverMock.On("GetBuildsForBuildType", "bcfg2", 1000).Times(1).Return(nil, errors.New("just and error to be ignored"))
				serverMock.On("GetBuildsForBuildType", "bcfg3", 1000).Times(1).Return(builds, nil)
				serverMock.On("GetBuildsForBuildType", "bcfg4", 1000).Times(1).Return(builds, nil)

				Convey("It should call GetBuildsForBuildType once for each build config", func() {
					tc.GetBuildHistory(&c)

					dbMock.AssertExpectations(t)
					serverMock.AssertExpectations(t)
				})
			})

			Convey("And there are builds", func() {
				b1 := teamcity.Build{BranchName: "dev", ID: 121}
				b2 := teamcity.Build{BranchName: "feat", ID: 122}
				b3 := teamcity.Build{BranchName: "dev", ID: 123}
				b4 := teamcity.Build{BranchName: "dev", ID: 124}
				b5 := teamcity.Build{BranchName: "feat", ID: 125}
				b6 := teamcity.Build{BranchName: "feat", ID: 126}

				serverMock.On("GetBuildsForBuildType", "bcfg1", 1000).Times(1).Return([]teamcity.Build{b1, b2, b3, b4}, nil)
				serverMock.On("GetBuildsForBuildType", "bcfg2", 1000).Times(1).Return([]teamcity.Build{b2, b3}, nil)
				serverMock.On("GetBuildsForBuildType", "bcfg3", 1000).Times(1).Return([]teamcity.Build{b1, b5, b3, b6, b2}, nil)
				serverMock.On("GetBuildsForBuildType", "bcfg4", 1000).Times(1).Return([]teamcity.Build{b1, b2, b6, b5}, nil)

				dbArray1 := []db.Branch{
					{Name: "dev", Builds: []db.Build{tc.BuildToDb(b1), tc.BuildToDb(b3), tc.BuildToDb(b4)}},
					{Name: "feat", Builds: []db.Build{tc.BuildToDb(b2)}},
				}
				dbArray2 := []db.Branch{
					{Name: "feat", Builds: []db.Build{tc.BuildToDb(b2)}},
					{Name: "dev", Builds: []db.Build{tc.BuildToDb(b3)}},
				}
				dbArray3 := []db.Branch{
					{Name: "dev", Builds: []db.Build{tc.BuildToDb(b1), tc.BuildToDb(b3)}},
					{Name: "feat", Builds: []db.Build{tc.BuildToDb(b5), tc.BuildToDb(b6), tc.BuildToDb(b2)}},
				}
				dbArray4 := []db.Branch{
					{Name: "dev", Builds: []db.Build{tc.BuildToDb(b1)}},
					{Name: "feat", Builds: []db.Build{tc.BuildToDb(b2), tc.BuildToDb(b6), tc.BuildToDb(b5)}},
				}

				dbMock.On("UpdateBuildTypeBuilds", "bcfg1", dbArray1).Times(1).Return(nil, nil)
				dbMock.On("UpdateBuildTypeBuilds", "bcfg2", dbArray2).Times(1).Return(nil, nil)
				dbMock.On("UpdateBuildTypeBuilds", "bcfg3", dbArray3).Times(1).Return(nil, errors.New("error to log and ignore"))
				dbMock.On("UpdateBuildTypeBuilds", "bcfg4", dbArray4).Times(1).Return(nil, nil)

				Convey("It should call GetBuildsForBuildType once for each build config And update the database with the builds", func() {

					tc.GetBuildHistory(&c)

					dbMock.AssertExpectations(t)
					serverMock.AssertExpectations(t)
				})
			})
		})
	})
}
