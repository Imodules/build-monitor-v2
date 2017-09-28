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
	"github.com/stretchr/testify/mock"
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

				var dbArray1 []db.Branch
				var dbArray2 []db.Branch
				var dbArray3 []db.Branch
				var dbArray4 []db.Branch

				dbMock.On("UpdateBuildTypeBuilds", "bcfg1", mock.Anything).Times(1).Return(nil, nil).Run(func(args mock.Arguments) {
					dbArray1 = args.Get(1).([]db.Branch)
				})
				dbMock.On("UpdateBuildTypeBuilds", "bcfg2", mock.Anything).Times(1).Return(nil, nil).Run(func(args mock.Arguments) {
					dbArray2 = args.Get(1).([]db.Branch)
				})
				dbMock.On("UpdateBuildTypeBuilds", "bcfg3", mock.Anything).Times(1).Return(nil, errors.New("error to log and ignore")).Run(func(args mock.Arguments) {
					dbArray3 = args.Get(1).([]db.Branch)
				})
				dbMock.On("UpdateBuildTypeBuilds", "bcfg4", mock.Anything).Times(1).Return(nil, nil).Run(func(args mock.Arguments) {
					dbArray4 = args.Get(1).([]db.Branch)
				})

				Convey("It should call GetBuildsForBuildType once for each build config And update the database with the builds", func() {

					tc.GetBuildHistory(&c)

					dbMock.AssertExpectations(t)
					serverMock.AssertExpectations(t)

					So(len(dbArray1), ShouldEqual, 2)
					So(dbArray1[0].Name, ShouldEqual, "dev")
					So(len(dbArray1[0].Builds), ShouldEqual, 3)

					So(dbArray1[1].Name, ShouldEqual, "feat")
					So(len(dbArray1[1].Builds), ShouldEqual, 1)

					So(len(dbArray2), ShouldEqual, 2)
					So(dbArray2[0].Name, ShouldEqual, "feat")
					So(len(dbArray2[0].Builds), ShouldEqual, 1)

					So(dbArray2[1].Name, ShouldEqual, "dev")
					So(len(dbArray2[1].Builds), ShouldEqual, 1)

					So(len(dbArray3), ShouldEqual, 2)
					So(dbArray3[0].Name, ShouldEqual, "dev")
					So(len(dbArray3[0].Builds), ShouldEqual, 2)

					So(dbArray3[1].Name, ShouldEqual, "feat")
					So(len(dbArray3[1].Builds), ShouldEqual, 3)

					So(len(dbArray4), ShouldEqual, 2)
					So(dbArray4[0].Name, ShouldEqual, "dev")
					So(len(dbArray4[0].Builds), ShouldEqual, 1)

					So(dbArray4[1].Name, ShouldEqual, "feat")
					So(len(dbArray4[1].Builds), ShouldEqual, 3)
				})
			})
		})
	})
}
