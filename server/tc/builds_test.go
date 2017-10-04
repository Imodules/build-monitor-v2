package tc_test

import (
	"testing"
	"time"

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
		tcMock := new(ITcClientMock)
		dbMock := new(IDbMock)

		c := tc.Server{
			Tc:  tcMock,
			Db:  dbMock,
			Log: log,
		}

		Convey("When GetRunningBuilds errors", func() {
			tcMock.On("GetRunningBuilds").Return([]teamcity.Build{}, errors.New("this shouldn't have happened"))

			lastBuilds := []teamcity.Build{{ID: 42}, {ID: 43}}
			newLastBuilds := tc.GetRunningBuilds(&c, lastBuilds)

			Convey("It should return the same build list we passed in", func() {
				tcMock.AssertExpectations(t)
				dbMock.AssertExpectations(t)

				So(newLastBuilds, ShouldResemble, lastBuilds)
			})
		})

		Convey("When GetRunningBuilds returns 0 builds", func() {
			tcMock.On("GetRunningBuilds").Return([]teamcity.Build{}, errors.New("this shouldn't have happened"))

			Convey("And no lastBuilds were passed in", func() {
				lastBuilds := []teamcity.Build{}

				newLastBuilds := tc.GetRunningBuilds(&c, lastBuilds)

				Convey("It should not do anything else", func() {
					tcMock.AssertExpectations(t)
					dbMock.AssertExpectations(t)

					So(len(newLastBuilds), ShouldEqual, 0)
				})
			})
		})

		Convey("When GetRunningBuilds returns builds", func() {
			runningBuilds := []teamcity.Build{
				{ID: 100, BuildTypeID: "bt100"}, // btErr
				{ID: 101, BuildTypeID: "bt101"}, // Still Processing
				{ID: 102, BuildTypeID: "bt102"}, // Ignore
				{ID: 104, BuildTypeID: "bt104"}, // New
				{ID: 105, BuildTypeID: "bt105"}, // New
			}

			lastBuilds := []teamcity.Build{
				{ID: 101, BuildTypeID: "bt101"}, // Still Processing
				{ID: 103, BuildTypeID: "bt103"}, // Completed
				{ID: 106, BuildTypeID: "bt106"}, // Completed (err FindBuildTypeById)
				{ID: 109, BuildTypeID: "bt109"}, // Completed (err GetBuildByID)
			}

			tcMock.On("GetRunningBuilds").Return(runningBuilds, nil)

			dbBt1 := db.BuildType{Id: "bt101", DashboardIds: []string{"abc", "123"}}
			dbBt2 := db.BuildType{Id: "bt102", DashboardIds: []string{}}
			dbBt4 := db.BuildType{Id: "bt104", DashboardIds: []string{"dash-1"}}
			dbBt5 := db.BuildType{Id: "bt105", DashboardIds: []string{"dash-1"}}

			dbMock.On("FindBuildTypeById", "bt100").Return(nil, errors.New("Something bad"))
			dbMock.On("FindBuildTypeById", "bt101").Return(&dbBt1, nil)
			dbMock.On("FindBuildTypeById", "bt102").Return(&dbBt2, nil)
			dbMock.On("FindBuildTypeById", "bt104").Return(&dbBt4, nil)
			dbMock.On("FindBuildTypeById", "bt105").Return(&dbBt5, nil)

			// Process last build
			dbBt3 := db.BuildType{Id: "bt103", DashboardIds: []string{"asdf"}}
			dbBt9 := db.BuildType{Id: "bt109", DashboardIds: []string{"asdf"}}
			dbMock.On("FindBuildTypeById", "bt103").Return(&dbBt3, nil)
			dbMock.On("FindBuildTypeById", "bt106").Return(nil, errors.New("this is an error"))
			dbMock.On("FindBuildTypeById", "bt109").Return(&dbBt9, nil)

			tcMock.On("GetBuildByID", 103).Return(teamcity.Build{ID: 103, Number: "api-103"}, nil)
			tcMock.On("GetBuildByID", 109).Return(teamcity.Build{}, errors.New("api failed"))

			Convey("And the ProcessRunningBuild succeeds", func() {
				oldProcessRunningBuild := tc.ProcessRunningBuild
				prbBuilds := []teamcity.Build{}
				prbBuildTypes := []db.BuildType{}
				tc.ProcessRunningBuild = func(c *tc.Server, b teamcity.Build, bt *db.BuildType) error {
					prbBuilds = append(prbBuilds, b)
					prbBuildTypes = append(prbBuildTypes, *bt)
					return nil
				}
				defer func() { tc.ProcessRunningBuild = oldProcessRunningBuild }()

				resultBuilds := tc.GetRunningBuilds(&c, lastBuilds)

				Convey("It should process all the builds in the builds list", func() {

					tcMock.AssertExpectations(t)
					dbMock.AssertExpectations(t)

					So(len(prbBuilds), ShouldEqual, 4)
					So(prbBuilds[0].ID, ShouldEqual, 101)
					So(prbBuilds[1].ID, ShouldEqual, 104)
					So(prbBuilds[2].ID, ShouldEqual, 105)
					So(prbBuilds[3].ID, ShouldEqual, 103)

					So(len(prbBuildTypes), ShouldEqual, 4)
					So(prbBuildTypes[0].Id, ShouldEqual, "bt101")
					So(prbBuildTypes[1].Id, ShouldEqual, "bt104")
					So(prbBuildTypes[2].Id, ShouldEqual, "bt105")
					So(prbBuildTypes[3].Id, ShouldEqual, "bt103")

					Convey("And it should finish any lastBuilds no longer in builds list", func() {
						Convey("And it should return useful builds", func() {
							So(len(resultBuilds), ShouldEqual, 3)
							So(resultBuilds[0].ID, ShouldEqual, 101)
							So(resultBuilds[1].ID, ShouldEqual, 104)
							So(resultBuilds[2].ID, ShouldEqual, 105)
						})
					})
				})
			})

			Convey("And the ProcessRunningBuild fails", func() {
				oldProcessRunningBuild := tc.ProcessRunningBuild
				prbBuilds := []teamcity.Build{}
				prbBuildTypes := []db.BuildType{}
				tc.ProcessRunningBuild = func(c *tc.Server, b teamcity.Build, bt *db.BuildType) error {
					prbBuilds = append(prbBuilds, b)
					prbBuildTypes = append(prbBuildTypes, *bt)
					return errors.New("i am failing")
				}
				defer func() { tc.ProcessRunningBuild = oldProcessRunningBuild }()

				resultBuilds := tc.GetRunningBuilds(&c, lastBuilds)

				Convey("It should process all the builds in the builds list", func() {

					tcMock.AssertExpectations(t)
					dbMock.AssertExpectations(t)

					So(len(prbBuilds), ShouldEqual, 4)
					So(prbBuilds[0].ID, ShouldEqual, 101)
					So(prbBuilds[1].ID, ShouldEqual, 104)
					So(prbBuilds[2].ID, ShouldEqual, 105)
					So(prbBuilds[3].ID, ShouldEqual, 103)

					So(len(prbBuildTypes), ShouldEqual, 4)
					So(prbBuildTypes[0].Id, ShouldEqual, "bt101")
					So(prbBuildTypes[1].Id, ShouldEqual, "bt104")
					So(prbBuildTypes[2].Id, ShouldEqual, "bt105")
					So(prbBuildTypes[3].Id, ShouldEqual, "bt103")

					Convey("And it should finish any lastBuilds no longer in builds list", func() {
						Convey("And it should return useful builds", func() {
							So(len(resultBuilds), ShouldEqual, 3)
							So(resultBuilds[0].ID, ShouldEqual, 101)
							So(resultBuilds[1].ID, ShouldEqual, 104)
							So(resultBuilds[2].ID, ShouldEqual, 105)
						})
					})
				})
			})
		})
	})
}

func TestServer_ProcessRunningBuild(t *testing.T) {
	Convey("Given a server", t, func() {
		log := logrus.WithField("test", "TestServer_ProcessRunningBuild")
		serverMock := new(ITcClientMock)
		dbMock := new(IDbMock)

		c := tc.Server{
			Tc:  serverMock,
			Db:  dbMock,
			Log: log,
		}

		Convey("When we have a build type without any branches", func() {

			tcBuild := teamcity.Build{
				ID:          801,
				BuildTypeID: "bt-id-801",
				BranchName:  "this is a b-name",
				Number:      "tc-build-number",
				Status:      teamcity.StatusRunning,
				StatusText:  "this will show up some place",
				Progress:    102,
				StartDate:   time.Unix(1507141495, 0),
				FinishDate:  time.Unix(1507141496, 0),
			}
			dbBuildType := db.BuildType{Id: "bt-id-801"}

			Convey("And the db update succeeds", func() {
				var branchesPassedToDb []db.Branch

				dbMock.On("UpdateBuildTypeBuilds", dbBuildType.Id, mock.AnythingOfType("[]db.Branch")).Return(nil, nil).Run(func(args mock.Arguments) {
					branchesPassedToDb = args.Get(1).([]db.Branch)
				})

				tc.ProcessRunningBuild(&c, tcBuild, &dbBuildType)

				Convey("It will create a new branch with the build And update the database", func() {
					dbMock.AssertExpectations(t)

					So(len(branchesPassedToDb), ShouldEqual, 1)
					So(branchesPassedToDb[0].Name, ShouldEqual, tcBuild.BranchName)

					So(len(branchesPassedToDb[0].Builds), ShouldEqual, 1)
					So(branchesPassedToDb[0].Builds[0].Id, ShouldEqual, tcBuild.ID)
					So(branchesPassedToDb[0].Builds[0].Number, ShouldEqual, tcBuild.Number)
					So(branchesPassedToDb[0].Builds[0].Status, ShouldEqual, tcBuild.Status)
					So(branchesPassedToDb[0].Builds[0].StatusText, ShouldEqual, tcBuild.StatusText)
					So(branchesPassedToDb[0].Builds[0].Progress, ShouldEqual, tcBuild.Progress)
					So(branchesPassedToDb[0].Builds[0].StartDate.Unix(), ShouldEqual, tcBuild.StartDate.Unix())
					So(branchesPassedToDb[0].Builds[0].FinishDate.Unix(), ShouldEqual, tcBuild.FinishDate.Unix())
				})
			})
		})

		Convey("When we have a build type that is already processing this build", func() {

			startDate := time.Now().Add(-1 * time.Minute)

			tcBuild := teamcity.Build{
				ID:          801,
				BuildTypeID: "bt-id-801",
				BranchName:  "this is a b-name",
				StartDate:   startDate,
			}

			dbBuildType := db.BuildType{
				Id: "bt-id-801",
				Branches: []db.Branch{
					{
						Name: "this is a b-name",
						Builds: []db.Build{
							{Id: 801},
							{Id: 799},
						},
					},
				},
			}

			Convey("And the db update succeeds", func() {
				var branchesPassedToDb []db.Branch

				dbMock.On("UpdateBuildTypeBuilds", dbBuildType.Id, mock.AnythingOfType("[]db.Branch")).Return(nil, nil).Run(func(args mock.Arguments) {
					branchesPassedToDb = args.Get(1).([]db.Branch)
				})

				now := time.Now()
				tc.ProcessRunningBuild(&c, tcBuild, &dbBuildType)

				Convey("It will update the build on the branch And update the database", func() {
					dbMock.AssertExpectations(t)

					So(len(branchesPassedToDb), ShouldEqual, 1)
					So(branchesPassedToDb[0].Name, ShouldEqual, tcBuild.BranchName)

					So(len(branchesPassedToDb[0].Builds), ShouldEqual, 2)
					So(branchesPassedToDb[0].Builds[0].Id, ShouldEqual, tcBuild.ID)
					So(branchesPassedToDb[0].Builds[0].StartDate.Unix(), ShouldEqual, startDate.Unix())
					So(branchesPassedToDb[0].Builds[0].FinishDate.Unix(), ShouldAlmostEqual, now.Unix())
					So(branchesPassedToDb[0].Builds[1].Id, ShouldEqual, 799)
				})
			})
		})

		Convey("When we have a build type that already has 12 builds and this is a new build", func() {
			tcBuild := teamcity.Build{
				ID:          801,
				BuildTypeID: "bt-id-801",
				BranchName:  "this is a b-name",
			}

			dbBuildType := db.BuildType{
				Id: "bt-id-801",
				Branches: []db.Branch{
					{
						Name: "this is a b-name",
						Builds: []db.Build{
							{Id: 800},
							{Id: 799},
							{Id: 798},
							{Id: 797},
							{Id: 796},
							{Id: 795},
							{Id: 794},
							{Id: 793},
							{Id: 792},
							{Id: 791},
							{Id: 790},
							{Id: 789},
							{Id: 788},
							{Id: 787},
							{Id: 786},
						},
					},
				},
			}

			Convey("And the db update succeeds", func() {
				var branchesPassedToDb []db.Branch

				dbMock.On("UpdateBuildTypeBuilds", dbBuildType.Id, mock.AnythingOfType("[]db.Branch")).Return(nil, nil).Run(func(args mock.Arguments) {
					branchesPassedToDb = args.Get(1).([]db.Branch)
				})

				tc.ProcessRunningBuild(&c, tcBuild, &dbBuildType)

				Convey("It will add the build to the beginning and remove any builds > 12 from the end And update the database", func() {
					dbMock.AssertExpectations(t)

					So(len(branchesPassedToDb), ShouldEqual, 1)
					So(branchesPassedToDb[0].Name, ShouldEqual, tcBuild.BranchName)

					So(len(branchesPassedToDb[0].Builds), ShouldEqual, 12)
					So(branchesPassedToDb[0].Builds[0].Id, ShouldEqual, tcBuild.ID)
					So(branchesPassedToDb[0].Builds[1].Id, ShouldEqual, 800)
					So(branchesPassedToDb[0].Builds[11].Id, ShouldEqual, 790)
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
