package db_test

import (
	"testing"
	"time"

	"build-monitor-v2/server/cfg"
	"build-monitor-v2/server/db"

	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2/bson"
)

func TestAppDb_UpsertBuildType(t *testing.T) {
	Convey("Given an appDb", t, func() {
		c := cfg.Config{PasswordSalt: "something here"}
		log := logrus.WithField("test", "TestAppDb_UpsertBuildType")

		appDb := db.Create(dbSession, &c, log, time.Now)

		buildType := db.BuildType{
			Id:        "Build Type Id 01",
			Name:      "Some build type for this",
			ProjectID: "Some project id here",
		}

		Convey("When a new buildType is inserted into the db", func() {
			result, err := appDb.UpsertBuildType(buildType)

			Convey("It should not error and return the object", func() {
				So(err, ShouldBeNil)
				So(result, ShouldNotBeNil)

				So(result.Id, ShouldEqual, buildType.Id)
				So(result.Name, ShouldEqual, buildType.Name)
				So(result.ProjectID, ShouldEqual, buildType.ProjectID)

				Convey("And it should be able to be found in the db", func() {
					var dbBuildType db.BuildType
					err := db.FindById(db.BuildTypes(appDb.Session), result.Id, &dbBuildType)

					So(err, ShouldBeNil)
					So(dbBuildType.Id, ShouldEqual, result.Id)
					So(dbBuildType.Name, ShouldEqual, result.Name)
					So(dbBuildType.ProjectID, ShouldEqual, result.ProjectID)
				})
			})
		})
	})
}

func TestAppDb_BuildTypeList(t *testing.T) {
	Convey("Given an appDb", t, func() {
		c := cfg.Config{PasswordSalt: "something here"}
		log := logrus.WithField("test", "TestAppDb_UpsertBuildType")

		db.BuildTypes(dbSession).RemoveAll(bson.M{})

		appDb := db.Create(dbSession, &c, log, time.Now)

		p1, _ := appDb.UpsertBuildType(db.BuildType{Id: "TestAppDb_BuildTypeList-p01", Name: "Ze End"})
		p2, _ := appDb.UpsertBuildType(db.BuildType{Id: "TestAppDb_BuildTypeList-p02", Name: "The End"})
		p3, _ := appDb.UpsertBuildType(db.BuildType{Id: "TestAppDb_BuildTypeList-p03", Name: "Deleted"})
		p4, _ := appDb.UpsertBuildType(db.BuildType{Id: "TestAppDb_BuildTypeList-p04", Name: "A good one"})

		err := appDb.DeleteBuildType(p3.Id)
		So(err, ShouldBeNil)

		Convey("When BuildTypeList is called", func() {
			projects, plErr := appDb.BuildTypeList()
			So(plErr, ShouldBeNil)

			Convey("It should return all non-deleted build types", func() {
				So(len(projects), ShouldEqual, 3)
				So(projects[0].Id, ShouldEqual, p4.Id)
				So(projects[1].Id, ShouldEqual, p2.Id)
				So(projects[2].Id, ShouldEqual, p1.Id)
			})
		})
	})
}
