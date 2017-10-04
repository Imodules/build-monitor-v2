package tc_test

import (
	"testing"

	"build-monitor-v2/server/tc"

	"build-monitor-v2/server/db"

	"errors"

	"github.com/pstuart2/go-teamcity"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestServer_GetRunningBuilds(t *testing.T) {
	Convey("Given a server", t, func() {
		log := logrus.WithField("test", "TestServer_GetBuildHistory")
		serverMock := new(ITcClientMock)
		dbMock := new(IDbMock)

		c := tc.Server{
			Tc:  serverMock,
			Db:  dbMock,
			Log: log,
		}

		Convey("When GetRunningBuilds errors", func() {
			serverMock.On("GetRunningBuilds").Return([]teamcity.Build{}, errors.New("this shouldn't have happened"))

			lastBuilds := []teamcity.Build{{ID: 42}, {ID: 43}}
			newLastBuilds := tc.GetRunningBuilds(&c, lastBuilds)

			Convey("It should return the same build list we passed in", func() {
				serverMock.AssertExpectations(t)
				dbMock.AssertExpectations(t)

				So(newLastBuilds, ShouldResemble, lastBuilds)
			})
		})

		Convey("When GetRunningBuilds returns 0 builds", func() {
			serverMock.On("GetRunningBuilds").Return([]teamcity.Build{}, errors.New("this shouldn't have happened"))

			Convey("And no lastBuilds were passed in", func() {
				lastBuilds := []teamcity.Build{}

				newLastBuilds := tc.GetRunningBuilds(&c, lastBuilds)

				Convey("It should not do anything else", func() {
					serverMock.AssertExpectations(t)
					dbMock.AssertExpectations(t)

					So(len(newLastBuilds), ShouldEqual, 0)
				})
			})
		})

		Convey("When GetRunningBuilds returns builds", func() {
			builds := []teamcity.Build{
				{ID: 100, BuildTypeID: "bt100"},                         // btErr
				{ID: 101, BuildTypeID: "bt101", BranchName: "branch-1"}, // Still Processing
				{ID: 102, BuildTypeID: "bt102"},                         // Ignore
				{ID: 104, BuildTypeID: "bt104"},                         // New
			}

			lastBuilds := []teamcity.Build{
				{ID: 101}, // Still Processing
				{ID: 103}, // Completed
			}

			serverMock.On("GetRunningBuilds").Return(builds, nil)

			dbBt1 := db.BuildType{Id: "bt1", DashboardIds: []string{"abc", "123"},
				Branches: []db.Branch{{Name: "branch-1"}},
			}

			dbBt2 := db.BuildType{Id: "bt1", DashboardIds: []string{}}
			//dbBt3 := db.BuildType{Id: "bt1", DashboardIds: []string{"asdf"}}
			dbBt4 := db.BuildType{Id: "bt1", DashboardIds: []string{"asdf"}}

			dbMock.On("FindBuildTypeById", "bt100").Return(nil, errors.New("Something bad"))
			dbMock.On("FindBuildTypeById", "bt101").Return(&dbBt1, nil)
			dbMock.On("FindBuildTypeById", "bt102").Return(&dbBt2, nil)
			//dbMock.On("FindBuildTypeById", "bt103").Return(&dbBt3, nil)
			dbMock.On("FindBuildTypeById", "bt104").Return(&dbBt4, nil)

			resultBuilds := tc.GetRunningBuilds(&c, lastBuilds)

			Convey("It should process all the builds in the builds list", func() {

				serverMock.AssertExpectations(t)
				dbMock.AssertExpectations(t)

				Convey("And it should finish any lastBuilds no longer in builds list", func() {
					Convey("And it should return useful builds", func() {
						So(len(resultBuilds), ShouldEqual, 2)
						So(resultBuilds[0].ID, ShouldEqual, 101)
						So(resultBuilds[1].ID, ShouldEqual, 104)
					})
				})
			})
		})
	})
}

func TestServer_GetBuildHistory(t *testing.T) {
	Convey("Given a server", t, func() {
		log := logrus.WithField("test", "TestServer_GetBuildHistory")
		serverMock := new(ITcClientMock)
		dbMock := new(IDbMock)

		c := tc.Server{
			Tc:  serverMock,
			Db:  dbMock,
			Log: log,
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
					So(len(getBranch("dev", dbArray1).Builds), ShouldEqual, 3)
					So(len(getBranch("feat", dbArray1).Builds), ShouldEqual, 1)

					So(len(dbArray2), ShouldEqual, 2)
					So(len(getBranch("dev", dbArray2).Builds), ShouldEqual, 1)
					So(len(getBranch("feat", dbArray2).Builds), ShouldEqual, 1)

					So(len(dbArray3), ShouldEqual, 2)
					So(len(getBranch("dev", dbArray3).Builds), ShouldEqual, 2)
					So(len(getBranch("feat", dbArray3).Builds), ShouldEqual, 3)

					So(len(dbArray4), ShouldEqual, 2)
					So(len(getBranch("dev", dbArray4).Builds), ShouldEqual, 1)
					So(len(getBranch("feat", dbArray4).Builds), ShouldEqual, 3)
				})
			})

			Convey("And there are > 12 builds", func() {
				b01 := teamcity.Build{BranchName: "dev", ID: 121}
				b02 := teamcity.Build{BranchName: "dev", ID: 122}
				b03 := teamcity.Build{BranchName: "dev", ID: 123}
				b04 := teamcity.Build{BranchName: "dev", ID: 124}
				b05 := teamcity.Build{BranchName: "dev", ID: 125}
				b06 := teamcity.Build{BranchName: "dev", ID: 126}
				b07 := teamcity.Build{BranchName: "dev", ID: 127}
				b08 := teamcity.Build{BranchName: "dev", ID: 128}
				b09 := teamcity.Build{BranchName: "dev", ID: 129}
				b10 := teamcity.Build{BranchName: "dev", ID: 130}
				b11 := teamcity.Build{BranchName: "dev", ID: 131}
				b12 := teamcity.Build{BranchName: "dev", ID: 132}
				b13 := teamcity.Build{BranchName: "dev", ID: 133}

				allBuilds := []teamcity.Build{b01, b02, b03, b04, b05, b06, b07, b08, b09, b10, b11, b12, b13}

				serverMock.On("GetBuildsForBuildType", "bcfg1", 1000).Times(1).Return(allBuilds, nil)
				serverMock.On("GetBuildsForBuildType", "bcfg2", 1000).Times(1).Return(allBuilds, nil)
				serverMock.On("GetBuildsForBuildType", "bcfg3", 1000).Times(1).Return(allBuilds, nil)
				serverMock.On("GetBuildsForBuildType", "bcfg4", 1000).Times(1).Return(allBuilds, nil)

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

				Convey("It should call GetBuildsForBuildType once for each build config And update the database with only 12 builds", func() {

					tc.GetBuildHistory(&c)

					dbMock.AssertExpectations(t)
					serverMock.AssertExpectations(t)

					So(len(dbArray1), ShouldEqual, 1)
					So(dbArray1[0].Name, ShouldEqual, "dev")
					So(len(dbArray1[0].Builds), ShouldEqual, 12)

					So(len(dbArray2), ShouldEqual, 1)
					So(dbArray2[0].Name, ShouldEqual, "dev")
					So(len(dbArray2[0].Builds), ShouldEqual, 12)

					So(len(dbArray3), ShouldEqual, 1)
					So(dbArray3[0].Name, ShouldEqual, "dev")
					So(len(dbArray3[0].Builds), ShouldEqual, 12)

					So(len(dbArray4), ShouldEqual, 1)
					So(dbArray4[0].Name, ShouldEqual, "dev")
					So(len(dbArray4[0].Builds), ShouldEqual, 12)
				})
			})
		})
	})
}

func getBranch(name string, branches []db.Branch) db.Branch {
	for _, b := range branches {
		if b.Name == name {
			return b
		}
	}

	return db.Branch{}
}
