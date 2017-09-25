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
			Id:           "Some random id -01",
			Name:         "Starting name that could change",
			OwnerId:      "this is an owner",
			BuildTypeIds: []string{"a1", "b2", "c3"},
		}

		Convey("When a new dashboard is inserted into the db", func() {
			result, err := appDb.UpsertDashboard(dashboard)

			Convey("It should not error and return the object", func() {
				So(err, ShouldBeNil)
				So(result, ShouldNotBeNil)

				So(result.Id, ShouldEqual, dashboard.Id)
				So(result.Name, ShouldEqual, dashboard.Name)
				So(result.OwnerId, ShouldEqual, dashboard.OwnerId)

				Convey("And it should be able to be found in the db", func() {
					var dbDashboard db.Dashboard
					err := db.FindById(db.Dashboards(appDb.Session), result.Id, &dbDashboard)

					So(err, ShouldBeNil)
					So(dbDashboard.Id, ShouldEqual, result.Id)
					So(dbDashboard.Name, ShouldEqual, result.Name)
					So(dbDashboard.OwnerId, ShouldEqual, result.OwnerId)

					So(len(dbDashboard.BuildTypeIds), ShouldEqual, 3)
				})
			})
		})

		Convey("When a dashboard already exists in the db", func() {
			dashboard.Name = "This is a new one!"
			dashboard.OwnerId = "This is where I belong"

			result, err := appDb.UpsertDashboard(dashboard)

			Convey("It should not error and return the updated object", func() {
				So(err, ShouldBeNil)
				So(result, ShouldNotBeNil)

				So(result.Id, ShouldEqual, dashboard.Id)
				So(result.Name, ShouldEqual, dashboard.Name)
				So(result.OwnerId, ShouldEqual, dashboard.OwnerId)

				Convey("And it should be able to be found in the db", func() {
					var dbDashboard db.Dashboard
					err := db.FindById(db.Dashboards(appDb.Session), result.Id, &dbDashboard)

					So(err, ShouldBeNil)
					So(dbDashboard.Id, ShouldEqual, result.Id)
					So(dbDashboard.Name, ShouldEqual, result.Name)
					So(dbDashboard.OwnerId, ShouldEqual, result.OwnerId)
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

		p1, _ := appDb.UpsertDashboard(db.Dashboard{Id: "TestAppDb_DashboardList-p01", Name: "Ze End", OwnerId: "ab1"})
		appDb.UpsertDashboard(db.Dashboard{Id: "TestAppDb_DashboardList-p02", Name: "The End", OwnerId: "ab2"})
		p3, _ := appDb.UpsertDashboard(db.Dashboard{Id: "TestAppDb_DashboardList-p03", Name: "Deleted", OwnerId: "ab1"})
		p4, _ := appDb.UpsertDashboard(db.Dashboard{Id: "TestAppDb_DashboardList-p04", Name: "A good one", OwnerId: "ab1"})

		err := appDb.DeleteDashboard(p3.Id)
		So(err, ShouldBeNil)

		Convey("When DashboardList is called", func() {
			dashboards, plErr := appDb.DashboardList("ab1")
			So(plErr, ShouldBeNil)

			Convey("It should return all non-deleted dashboards", func() {
				So(len(dashboards), ShouldEqual, 2)
				So(dashboards[0].Id, ShouldEqual, p4.Id)
				So(dashboards[1].Id, ShouldEqual, p1.Id)
			})
		})
	})
}
