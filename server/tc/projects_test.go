package tc_test

import (
	"testing"

	"build-monitor-v2/server/tc"

	"time"

	"github.com/kapitanov/go-teamcity"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestServer_RefreshProjects(t *testing.T) {
	Convey("Given a server", t, func() {
		log := logrus.WithField("test", "TestServer_Start_Shutdown")
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

			Convey("It should not call the db", func() {
				c.RefreshProjects()
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

			db1 := tc.ProjectToDb(projects[1])
			db2 := tc.ProjectToDb(projects[2])
			db3 := tc.ProjectToDb(projects[3])
			db4 := tc.ProjectToDb(projects[4])

			dbMock.On("UpsertProject", db1).Return(&db1, nil)
			dbMock.On("UpsertProject", db2).Return(&db2, nil)
			dbMock.On("UpsertProject", db3).Return(&db3, nil)
			dbMock.On("UpsertProject", db4).Return(&db4, nil)

			Convey("It should call the db for each project that is not _Root", func() {
				c.RefreshProjects()

				dbMock.AssertExpectations(t)
			})
		})
	})
}
