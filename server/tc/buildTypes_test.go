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

func TestServer_RefreshBuildTypes(t *testing.T) {
	Convey("Given a server", t, func() {
		log := logrus.WithField("test", "TestServer_RefreshBuildTypes")
		serverMock := new(ITcClientMock)
		dbMock := new(IDbMock)

		c := tc.Server{
			Tc:                  serverMock,
			Db:                  dbMock,
			Log:                 log,
			ProjectPollInterval: time.Millisecond * 500,
		}

		Convey("When there are no projects", func() {
			projects := []db.Project{}

			dbMock.On("ProjectList").Return(projects, nil)

			Convey("It should not query team city or update the db", func() {
				tc.RefreshBuildTypes(&c)

				serverMock.AssertExpectations(t)
				dbMock.AssertExpectations(t)
			})
		})

		Convey("When we fail to get the project map", func() {
			expectedError := errors.New("this was expected")
			dbMock.On("ProjectList").Return(nil, expectedError)

			Convey("It should not query team city or update the db and return the error", func() {
				err := tc.RefreshBuildTypes(&c)

				serverMock.AssertExpectations(t)
				dbMock.AssertExpectations(t)

				So(err, ShouldEqual, expectedError)
			})
		})

		Convey("When there are no build types", func() {
			projects := []db.Project{
				{Id: "something neat 1"},
			}

			dbMock.On("ProjectList").Return(projects, nil)

			buildTypes := []teamcity.BuildType{}

			serverMock.On("GetBuildTypes").Times(1).Return(buildTypes, nil)

			dbBuildTypes := []db.BuildType{}
			dbMock.On("BuildTypeList").Return(dbBuildTypes, nil)

			Convey("It should not update the db", func() {
				tc.RefreshBuildTypes(&c)

				serverMock.AssertExpectations(t)
				dbMock.AssertExpectations(t)
			})
		})

		Convey("When there are buildTypes and projects", func() {
			projects := []db.Project{
				{Id: "p1"},
				{Id: "p3"},
			}

			dbMock.On("ProjectList").Return(projects, nil)

			buildTypes := []teamcity.BuildType{
				{ID: "bt1", Name: "1 BuildType", Description: "Something here _", ProjectID: "p1"},
				{ID: "bt2", Name: "2 BuildType", Description: "Something here 1", ProjectID: "p3"},
				{ID: "bt3", Name: "3 BuildType", Description: "Something here 2", ProjectID: "p1"},
				{ID: "bt4", Name: "4 BuildType", Description: "Something here 3", ProjectID: "p3"},
				{ID: "bt5", Name: "5 BuildType", Description: "Something here 3", ProjectID: "not-used"},
			}

			serverMock.On("GetBuildTypes").Times(1).Return(buildTypes, nil)

			Convey("And the get BuildTypeList is successful", func() {
				dbBuildTypes := []db.BuildType{}
				dbMock.On("BuildTypeList").Return(dbBuildTypes, nil)

				db1 := tc.BuildTypeToDb(buildTypes[0])
				db2 := tc.BuildTypeToDb(buildTypes[1])
				db3 := tc.BuildTypeToDb(buildTypes[2])
				db4 := tc.BuildTypeToDb(buildTypes[3])
				//db5 := tc.BuildTypeToDb(buildTypes[4])

				Convey("And the upsert is successful", func() {
					dbMock.On("UpsertBuildType", db1).Return(&db1, nil)
					dbMock.On("UpsertBuildType", db2).Return(&db2, nil)
					dbMock.On("UpsertBuildType", db3).Return(&db3, nil)
					dbMock.On("UpsertBuildType", db4).Return(&db4, nil)

					Convey("It should call the db for each project that is not _Root", func() {
						tc.RefreshBuildTypes(&c)

						serverMock.AssertExpectations(t)
						dbMock.AssertExpectations(t)
					})
				})

				Convey("And the upsert fails", func() {
					someError := errors.New("some error that will just get logged")

					dbMock.On("UpsertBuildType", db1).Return(&db1, nil)
					dbMock.On("UpsertBuildType", db2).Return(nil, someError)
					dbMock.On("UpsertBuildType", db3).Return(nil, someError)
					dbMock.On("UpsertBuildType", db4).Return(&db4, nil)

					Convey("It should call the db for each project that is not _Root", func() {
						tc.RefreshBuildTypes(&c)

						serverMock.AssertExpectations(t)
						dbMock.AssertExpectations(t)
					})
				})
			})

			Convey("And the get BuildTypeList fails", func() {
				expectedError := errors.New("f'd to get project list")

				dbMock.On("BuildTypeList").Return(nil, expectedError)

				Convey("It should return the error", func() {
					err := tc.RefreshBuildTypes(&c)
					So(err, ShouldEqual, expectedError)
				})
			})
		})

		Convey("When there are buildTypes in the database that are no longer in Tc", func() {
			projects := []db.Project{
				{Id: "p1"},
				{Id: "p3"},
			}

			dbMock.On("ProjectList").Return(projects, nil)

			buildTypes := []teamcity.BuildType{
				{ID: "bt1", Name: "1 BuildType", Description: "Something here _", ProjectID: "p1"},
				{ID: "bt2", Name: "2 BuildType", Description: "Something here 1", ProjectID: "p3"},
				{ID: "bt3", Name: "3 BuildType", Description: "Something here 2", ProjectID: "p1"},
				{ID: "bt4", Name: "4 BuildType", Description: "Something here 3", ProjectID: "p3"},
			}

			serverMock.On("GetBuildTypes").Times(1).Return(buildTypes, nil)

			db1 := tc.BuildTypeToDb(buildTypes[0])
			db2 := tc.BuildTypeToDb(buildTypes[1])
			db3 := tc.BuildTypeToDb(buildTypes[2])
			db4 := tc.BuildTypeToDb(buildTypes[3])
			db5 := tc.BuildTypeToDb(teamcity.BuildType{ID: "p5", Name: "5 BuildType", Description: "Something here 5", ProjectID: "p1"})
			db6 := tc.BuildTypeToDb(teamcity.BuildType{ID: "p6", Name: "6 BuildType", Description: "Something here 6", ProjectID: "p3"})

			dbBuildTypes := []db.BuildType{db1, db2, db3, db4, db5, db6}
			dbMock.On("BuildTypeList").Return(dbBuildTypes, nil)
			dbMock.On("DeleteBuildType", db5.Id).Return(nil)
			dbMock.On("DeleteBuildType", db6.Id).Return(nil)

			dbMock.On("UpsertBuildType", db1).Return(&db1, nil)
			dbMock.On("UpsertBuildType", db2).Return(&db2, nil)
			dbMock.On("UpsertBuildType", db3).Return(&db3, nil)
			dbMock.On("UpsertBuildType", db4).Return(&db4, nil)

			Convey("It should call the db for each buildType", func() {
				tc.RefreshBuildTypes(&c)

				serverMock.AssertExpectations(t)
				dbMock.AssertExpectations(t)
			})
		})

		Convey("When we fail to get the tc buildTypes", func() {
			projects := []db.Project{
				{Id: "p1"},
				{Id: "p3"},
			}

			dbMock.On("ProjectList").Return(projects, nil)

			expectedErr := errors.New("this was expected")
			serverMock.On("GetBuildTypes").Times(1).Return(nil, expectedErr)

			Convey("It should call the db for each project that is not _Root", func() {
				err := tc.RefreshBuildTypes(&c)

				So(err, ShouldEqual, expectedErr)

				serverMock.AssertExpectations(t)
				dbMock.AssertExpectations(t)
			})
		})
	})
}
