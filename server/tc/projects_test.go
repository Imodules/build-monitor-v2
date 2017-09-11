package tc_test

import (
	"testing"

	"build-monitor-v2/server/tc"

	"time"

	"build-monitor-v2/server/db"
	"errors"
	"github.com/kapitanov/go-teamcity"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestServer_RefreshProjects(t *testing.T) {
	Convey("Given a server", t, func() {
		log := logrus.WithField("test", "TestServer_RefreshProjects")
		serverMock := new(ITcClientMock)
		dbMock := new(IDbMock)

		c := tc.Server{
			Tc:                  serverMock,
			Db:                  dbMock,
			Log:                 log,
			ProjectPollInterval: time.Millisecond * 500,
		}

		Convey("When there are no projects", func() {
			projects := []teamcity.Project{}

			serverMock.On("GetProjects").Times(1).Return(projects, nil)

			dbProjects := []db.Project{}
			dbMock.On("ProjectList").Return(dbProjects, nil)

			Convey("It should not update the db", func() {
				tc.RefreshProjects(&c)

				serverMock.AssertExpectations(t)
				dbMock.AssertExpectations(t)
			})
		})

		Convey("When there are projects", func() {
			projects := []teamcity.Project{
				{ID: "_Root", Name: "Root Project", Description: "I am _Root"},
				{ID: "p1", Name: "1 Project", Description: "Something here _", ParentProjectID: "_Root"},
				{ID: "p2", Name: "2 Project", Description: "Something here 1", ParentProjectID: "_Root"},
				{ID: "p3", Name: "3 Project", Description: "Something here 2", ParentProjectID: "p1"},
				{ID: "p4", Name: "4 Project", Description: "Something here 3", ParentProjectID: "p3"},
			}

			serverMock.On("GetProjects").Times(1).Return(projects, nil)

			Convey("And the get ProjectList is successful", func() {
				dbProjects := []db.Project{}
				dbMock.On("ProjectList").Return(dbProjects, nil)

				db1 := tc.ProjectToDb(projects[1])
				db2 := tc.ProjectToDb(projects[2])
				db3 := tc.ProjectToDb(projects[3])
				db4 := tc.ProjectToDb(projects[4])

				Convey("And the upsert is successful", func() {
					dbMock.On("UpsertProject", db1).Return(&db1, nil)
					dbMock.On("UpsertProject", db2).Return(&db2, nil)
					dbMock.On("UpsertProject", db3).Return(&db3, nil)
					dbMock.On("UpsertProject", db4).Return(&db4, nil)

					Convey("It should call the db for each project that is not _Root", func() {
						tc.RefreshProjects(&c)

						serverMock.AssertExpectations(t)
						dbMock.AssertExpectations(t)
					})
				})

				Convey("And the upsert fails", func() {
					someError := errors.New("some error that will just get logged")

					dbMock.On("UpsertProject", db1).Return(&db1, nil)
					dbMock.On("UpsertProject", db2).Return(nil, someError)
					dbMock.On("UpsertProject", db3).Return(nil, someError)
					dbMock.On("UpsertProject", db4).Return(&db4, nil)

					Convey("It should call the db for each project that is not _Root", func() {
						tc.RefreshProjects(&c)

						serverMock.AssertExpectations(t)
						dbMock.AssertExpectations(t)
					})
				})
			})

			Convey("And the get ProjectList fails", func() {
				expectedError := errors.New("f'd to get project list")

				dbMock.On("ProjectList").Return(nil, expectedError)

				Convey("It should return the error", func() {
					err := tc.RefreshProjects(&c)
					So(err, ShouldEqual, expectedError)
				})
			})
		})

		Convey("When there are projects in the database that are no longer in Tc", func() {
			projects := []teamcity.Project{
				{ID: "_Root", Name: "Root Project", Description: "I am _Root"},
				{ID: "p1", Name: "1 Project", Description: "Something here _", ParentProjectID: "_Root"},
				{ID: "p2", Name: "2 Project", Description: "Something here 1", ParentProjectID: "_Root"},
				{ID: "p3", Name: "3 Project", Description: "Something here 2", ParentProjectID: "p1"},
				{ID: "p4", Name: "4 Project", Description: "Something here 3", ParentProjectID: "p3"},
			}

			serverMock.On("GetProjects").Times(1).Return(projects, nil)

			db1 := tc.ProjectToDb(projects[1])
			db2 := tc.ProjectToDb(projects[2])
			db3 := tc.ProjectToDb(projects[3])
			db4 := tc.ProjectToDb(projects[4])
			db5 := tc.ProjectToDb(teamcity.Project{ID: "p5", Name: "5 Project", Description: "Something here 5", ParentProjectID: "p1"})
			db6 := tc.ProjectToDb(teamcity.Project{ID: "p6", Name: "6 Project", Description: "Something here 6", ParentProjectID: "p3"})

			dbProjects := []db.Project{db1, db2, db3, db4, db5, db6}
			dbMock.On("ProjectList").Return(dbProjects, nil)
			dbMock.On("DeleteProject", db5.Id).Return(nil)
			dbMock.On("DeleteProject", db6.Id).Return(nil)

			dbMock.On("UpsertProject", db1).Return(&db1, nil)
			dbMock.On("UpsertProject", db2).Return(&db2, nil)
			dbMock.On("UpsertProject", db3).Return(&db3, nil)
			dbMock.On("UpsertProject", db4).Return(&db4, nil)

			Convey("It should call the db for each project that is not _Root", func() {
				tc.RefreshProjects(&c)

				serverMock.AssertExpectations(t)
				dbMock.AssertExpectations(t)
			})
		})

		Convey("When we fail to get the tc projects", func() {
			expectedErr := errors.New("this was expected")
			serverMock.On("GetProjects").Times(1).Return(nil, expectedErr)

			Convey("It should call the db for each project that is not _Root", func() {
				err := tc.RefreshProjects(&c)

				So(err, ShouldEqual, expectedErr)

				serverMock.AssertExpectations(t)
				dbMock.AssertExpectations(t)
			})
		})
	})
}
