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

func TestAppDb_UpsertDashboard(t *testing.T) {
	Convey("Given an AppDb", t, func() {
		c := cfg.Config{PasswordSalt: "something here"}
		log := logrus.WithField("test", "TestAppDb_UpsertDashboard")

		appDb := db.Create(dbSession, &c, log, time.Now)

		dashboard := db.Dashboard{
			Id:          "Some random id -01",
			Name:        "Starting name that could change",
			Owner:       db.Owner{Id: bson.NewObjectId(), Username: "cool me"},
			ColumnCount: 4,
			BuildConfigs: []db.BuildConfig{
				{Id: "a1", Abbreviation: "cool 1"},
				{Id: "b2", Abbreviation: "cool 2"},
				{Id: "c3", Abbreviation: "cool 3"},
			},
		}

		Convey("When a new dashboard is inserted into the db", func() {
			result, err := appDb.UpsertDashboard(dashboard)

			Convey("It should not error and return the object", func() {
				So(err, ShouldBeNil)
				So(result, ShouldNotBeNil)

				So(result.Id, ShouldEqual, dashboard.Id)
				So(result.Name, ShouldEqual, dashboard.Name)
				So(result.Owner.Id.Hex(), ShouldEqual, dashboard.Owner.Id.Hex())
				So(result.Owner.Username, ShouldEqual, dashboard.Owner.Username)
				So(len(result.BuildConfigs), ShouldEqual, 3)

				Convey("And it should be able to be found in the db", func() {
					var dbDashboard db.Dashboard
					err := db.FindById(db.Dashboards(appDb.Session), result.Id, &dbDashboard)

					So(err, ShouldBeNil)
					So(dbDashboard.Id, ShouldEqual, result.Id)
					So(dbDashboard.Name, ShouldEqual, result.Name)
					So(dbDashboard.ColumnCount, ShouldEqual, result.ColumnCount)
					So(dbDashboard.Owner.Id.Hex(), ShouldEqual, result.Owner.Id.Hex())
					So(dbDashboard.Owner.Username, ShouldEqual, result.Owner.Username)

					So(len(dbDashboard.BuildConfigs), ShouldEqual, 3)
					So(dbDashboard.BuildConfigs[1].Id, ShouldEqual, "b2")
					So(dbDashboard.BuildConfigs[1].Abbreviation, ShouldEqual, "cool 2")
				})
			})
		})

		Convey("When a dashboard already exists in the db", func() {
			dashboard.Name = "This is a new one!"
			dashboard.BuildConfigs[0].Id = "new 1"
			dashboard.BuildConfigs[1].Abbreviation = "Changing it"

			result, err := appDb.UpsertDashboard(dashboard)

			Convey("It should not error and return the updated object", func() {
				So(err, ShouldBeNil)
				So(result, ShouldNotBeNil)

				So(result.Id, ShouldEqual, dashboard.Id)
				So(result.Name, ShouldEqual, dashboard.Name)

				Convey("And it should be able to be found in the db", func() {
					var dbDashboard db.Dashboard
					err := db.FindById(db.Dashboards(appDb.Session), result.Id, &dbDashboard)

					So(err, ShouldBeNil)
					So(dbDashboard.Id, ShouldEqual, result.Id)
					So(dbDashboard.Name, ShouldEqual, result.Name)
					So(len(dashboard.BuildConfigs), ShouldEqual, 3)
					So(dbDashboard.BuildConfigs[0].Id, ShouldEqual, "new 1")
					So(dbDashboard.BuildConfigs[1].Abbreviation, ShouldEqual, "Changing it")
				})
			})
		})

		Convey("When dashboard is deleted in the db", func() {
			delErr := appDb.DeleteDashboard(dashboard.Id)
			So(delErr, ShouldBeNil)

			dashboard.Name = "again!"

			result, err := appDb.UpsertDashboard(dashboard)

			Convey("It should not error and return the updated object", func() {
				So(err, ShouldBeNil)
				So(result, ShouldNotBeNil)

				So(result.Id, ShouldEqual, dashboard.Id)
				So(result.Name, ShouldEqual, dashboard.Name)

				Convey("And it should remove the deleted flag and be able to be found in the db", func() {
					var dbDashboard db.Dashboard
					err := db.FindById(db.Dashboards(appDb.Session), result.Id, &dbDashboard)

					So(err, ShouldBeNil)
					So(dbDashboard.Id, ShouldEqual, result.Id)
					So(dbDashboard.Name, ShouldEqual, result.Name)
				})
			})
		})
	})
}

func TestAppDb_DashboardList(t *testing.T) {
	Convey("Given an appDb", t, func() {
		c := cfg.Config{PasswordSalt: "something here"}
		log := logrus.WithField("test", "TestAppDb_UpsertDashboard")

		db.Dashboards(dbSession).RemoveAll(bson.M{})

		appDb := db.Create(dbSession, &c, log, time.Now)

		p1, _ := appDb.UpsertDashboard(db.Dashboard{Id: "TestAppDb_DashboardList-p01", Name: "Ze End"})
		p2, _ := appDb.UpsertDashboard(db.Dashboard{Id: "TestAppDb_DashboardList-p02", Name: "The End"})
		p3, _ := appDb.UpsertDashboard(db.Dashboard{Id: "TestAppDb_DashboardList-p03", Name: "Deleted"})
		p4, _ := appDb.UpsertDashboard(db.Dashboard{Id: "TestAppDb_DashboardList-p04", Name: "A good one"})

		err := appDb.DeleteDashboard(p3.Id)
		So(err, ShouldBeNil)

		Convey("When DashboardList is called", func() {
			dashboards, plErr := appDb.DashboardList()
			So(plErr, ShouldBeNil)

			Convey("It should return all non-deleted dashboards", func() {
				So(len(dashboards), ShouldEqual, 3)
				So(dashboards[0].Id, ShouldEqual, p4.Id)
				So(dashboards[1].Id, ShouldEqual, p2.Id)
				So(dashboards[2].Id, ShouldEqual, p1.Id)
			})
		})
	})
}
