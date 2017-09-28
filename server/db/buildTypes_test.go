package db_test

import (
	"testing"
	"time"

	"build-monitor-v2/server/cfg"
	"build-monitor-v2/server/db"

	"github.com/pstuart2/go-teamcity"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2/bson"
)

func TestAppDb_UpsertBuildType(t *testing.T) {
	Convey("Given an appDb", t, func() {
		c := cfg.Config{PasswordSalt: "something here"}
		log := logrus.WithField("test", "TestAppDb_UpsertBuildType")

		appDb := db.Create(dbSession, &c, log, time.Now)

		Convey("When a new buildType is inserted into the db", func() {
			buildType := db.BuildType{
				Id:        "Build Type Id 01",
				Name:      "Some build type for this",
				ProjectID: "Some project id here",
			}

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

		Convey("When an existing buildType with builds is updated in the db", func() {
			buildType := db.BuildType{
				Id:        "Build Type Id 02",
				Name:      "Some build type for this 2",
				ProjectID: "Some project id here 2",
			}

			bt, _ := appDb.UpsertBuildType(buildType)

			builds := []db.Build{
				{Id: 908, Number: "BT-51", Status: teamcity.StatusSuccess, StatusText: "this was the last", Progress: 100, BranchName: "master"},
				{Id: 905, Number: "BT-41", Status: teamcity.StatusFailure, StatusText: "dddd", Progress: 88, BranchName: "master"},
				{Id: 901, Number: "BT-31", Status: teamcity.StatusSuccess, StatusText: "ffff", Progress: 55, BranchName: "master"},
			}

			appDb.UpdateBuildTypeBuilds(bt.Id, builds)

			Convey("It should not remove the builds from the db", func() {
				dbBuildType, _ := appDb.FindBuildTypeById(bt.Id)

				So(len(dbBuildType.Builds), ShouldEqual, 3)

				So(dbBuildType.Builds[0].Id, ShouldEqual, 908)
				So(dbBuildType.Builds[1].Id, ShouldEqual, 905)
				So(dbBuildType.Builds[2].Id, ShouldEqual, 901)
			})
		})
	})
}

func TestAppDb_UpdateBuildTypeBuilds(t *testing.T) {
	Convey("Given an appDb", t, func() {
		c := cfg.Config{PasswordSalt: "something here"}
		log := logrus.WithField("test", "TestAppDb_UpsertBuildType")

		appDb := db.Create(dbSession, &c, log, time.Now)

		buildType := db.BuildType{
			Id:        "Build Type Id 01",
			Name:      "Some build type for this",
			ProjectID: "Some project id here",
		}

		Convey("With a valid build type", func() {
			bt, btError := appDb.UpsertBuildType(buildType)
			So(btError, ShouldBeNil)

			Convey("When we add builds to the build type", func() {

				builds := []db.Build{
					{Id: 908, Number: "BT-51", Status: teamcity.StatusSuccess, StatusText: "this was the last", Progress: 100, BranchName: "master"},
					{Id: 905, Number: "BT-41", Status: teamcity.StatusFailure, StatusText: "dddd", Progress: 88, BranchName: "master"},
					{Id: 901, Number: "BT-31", Status: teamcity.StatusSuccess, StatusText: "ffff", Progress: 55, BranchName: "master"},
				}

				newBt, err := appDb.UpdateBuildTypeBuilds(bt.Id, builds)

				Convey("It should not error", func() {
					So(err, ShouldBeNil)

					Convey("And return the updated build type", func() {
						So(len(newBt.Builds), ShouldEqual, 3)

						So(newBt.Builds[0].Id, ShouldEqual, 908)
						So(newBt.Builds[1].Id, ShouldEqual, 905)
						So(newBt.Builds[2].Id, ShouldEqual, 901)

						Convey("And we should be able to query them from the db", func() {
							dbBuildType, dbErr := appDb.FindBuildTypeById(bt.Id)

							So(dbErr, ShouldBeNil)
							So(dbBuildType.Id, ShouldEqual, bt.Id)
							So(len(dbBuildType.Builds), ShouldEqual, 3)

							So(dbBuildType.Builds[0].Id, ShouldEqual, 908)
							So(dbBuildType.Builds[1].Id, ShouldEqual, 905)
							So(dbBuildType.Builds[2].Id, ShouldEqual, 901)

							So(dbBuildType.Builds[0].Number, ShouldEqual, "BT-51")
							So(dbBuildType.Builds[1].Number, ShouldEqual, "BT-41")
							So(dbBuildType.Builds[2].Number, ShouldEqual, "BT-31")
						})
					})
				})
			})
		})
	})
}

func TestAppDb_AddRemoveDashboardFromBuildTypes(t *testing.T) {
	Convey("Given an appDb", t, func() {
		c := cfg.Config{PasswordSalt: "something here"}
		log := logrus.WithField("test", "TestAppDb_UpsertBuildType")

		appDb := db.Create(dbSession, &c, log, time.Now)

		Convey("With multiple build types with dashboard ids", func() {
			buildType1 := db.BuildType{
				Id: "Build Type Id 01",
			}
			bt1, _ := appDb.UpsertBuildType(buildType1)

			buildType2 := db.BuildType{
				Id: "Build Type Id 02",
			}
			bt2, _ := appDb.UpsertBuildType(buildType2)

			buildType3 := db.BuildType{
				Id: "Build Type Id 03",
			}
			bt3, _ := appDb.UpsertBuildType(buildType3)

			appDb.AddDashboardToBuildTypes([]string{bt1.Id, bt3.Id}, "dash-01")
			appDb.AddDashboardToBuildTypes([]string{bt1.Id, bt2.Id, bt3.Id}, "dash-02")
			appDb.AddDashboardToBuildTypes([]string{bt2.Id, bt3.Id}, "dash-03")

			Convey("When I remove a dashboard from build types", func() {
				err := appDb.RemoveDashboardFromBuildTypes("dash-02")
				So(err, ShouldBeNil)

				Convey("It should remove them from all", func() {
					dash2Results, err := appDb.DashboardBuildTypeList("dash-02")
					So(err, ShouldBeNil)
					So(len(dash2Results), ShouldEqual, 0)

					Convey("But only that id", func() {
						dash1Results, err := appDb.DashboardBuildTypeList("dash-01")
						So(err, ShouldBeNil)

						So(len(dash1Results), ShouldEqual, 2)
						So(dash1Results[0].Id, ShouldEqual, bt1.Id)
						So(dash1Results[1].Id, ShouldEqual, bt3.Id)

						dash3Results, err := appDb.DashboardBuildTypeList("dash-03")
						So(err, ShouldBeNil)

						So(len(dash3Results), ShouldEqual, 2)
						So(dash3Results[0].Id, ShouldEqual, bt2.Id)
						So(dash3Results[1].Id, ShouldEqual, bt3.Id)
					})
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
