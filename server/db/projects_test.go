package db_test

import (
	"testing"
	"time"

	"build-monitor-v2/server/cfg"
	"build-monitor-v2/server/db"

	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAppDb_UpsertProject(t *testing.T) {
	Convey("Given an AppDb", t, func() {
		c := cfg.Config{PasswordSalt: "something here"}
		log := logrus.WithField("test", "TestAppDb_UpsertProject")

		appDb := db.Create(dbSession, &c, log, time.Now)

		project := db.Project{
			Id:   "This is unique id 001",
			Name: "Starting name that could change",
		}

		Convey("When a new project is inserted into the db", func() {
			result, err := appDb.UpsertProject(project)

			Convey("It should not error and return the object", func() {
				So(err, ShouldBeNil)
				So(result, ShouldNotBeNil)

				So(result.Id, ShouldEqual, project.Id)
				So(result.Name, ShouldEqual, project.Name)

				Convey("And it should be able to be found in the db", func() {
					var dbProject db.Project
					err := db.FindById(db.Projects(appDb.Session), result.Id, &dbProject)

					So(err, ShouldBeNil)
					So(dbProject.Id, ShouldEqual, result.Id)
					So(dbProject.Name, ShouldEqual, result.Name)
				})
			})
		})

		Convey("When a project already exists in the db", func() {
			project.Name = "This is a new one!"
			project.Description = "For this is true"
			project.ParentProjectID = "This is where I belong"

			result, err := appDb.UpsertProject(project)

			Convey("It should not error and return the updated object", func() {
				So(err, ShouldBeNil)
				So(result, ShouldNotBeNil)

				So(result.Id, ShouldEqual, project.Id)
				So(result.Name, ShouldEqual, project.Name)

				Convey("And it should be able to be found in the db", func() {
					var dbProject db.Project
					err := db.FindById(db.Projects(appDb.Session), result.Id, &dbProject)

					So(err, ShouldBeNil)
					So(dbProject.Id, ShouldEqual, result.Id)
					So(dbProject.Name, ShouldEqual, result.Name)
				})
			})
		})
	})
}
