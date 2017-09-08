package db_test

import (
	"os"
	"testing"

	"build-monitor-v2/server/cfg"

	"build-monitor-v2/server/db"

	"fmt"

	"time"

	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2"
)

var dbSession *mgo.Session
var userCounter = 0

func resetDb() {
	dbSession.DB("").DropDatabase()
}

func getNewUser(appDb *db.AppDb, name string) *db.User {
	userCounter++
	username := fmt.Sprintf("%s-%d", name, userCounter)

	user, _ := appDb.CreateUser(username, username+"@fwe.com", username)
	return user
}

func TestMain(m *testing.M) {
	session, _ := mgo.Dial("mongodb://localhost/build-monitor-v2-test")
	defer session.Close()

	dbSession = session
	resetDb()

	log := logrus.WithField("test", "data/TestMain")
	if err := db.Ensure(session, log); err != nil {
		log.Fatalf("Failed to ensure the appDb: %v", err)
	}

	code := m.Run()

	os.Exit(code)
}

func TestCreate(t *testing.T) {
	Convey("Given a session and a config", t, func() {
		log := logrus.WithField("test", "dbLog")
		c := cfg.Config{PasswordSalt: "you know it!"}

		Convey("Then Create should return an AppDb object", func() {
			appDb := db.Create(dbSession, &c, log, time.Now)

			So(appDb.Session, ShouldEqual, dbSession)
			So(appDb.PasswordSalt, ShouldEqual, c.PasswordSalt)
			So(appDb.Log, ShouldEqual, log)

		})
	})
}

func TestAppDb_Delete(t *testing.T) {
	Convey("Given an AppDb with a record", t, func() {
		c := cfg.Config{PasswordSalt: "something here"}
		log := logrus.WithField("test", "TestAppDb_UpsertProject")

		appDb := db.Create(dbSession, &c, log, time.Now)

		project := db.Project{
			Id:   "This is one I will delete",
			Name: "Starting name that could change",
		}

		result, err := appDb.UpsertProject(project)
		So(err, ShouldBeNil)
		So(result, ShouldNotBeNil)

		Convey("When the record is deleted", func() {
			err := appDb.Delete(db.Projects(appDb.Session), result.Id)

			Convey("It should not error", func() {
				So(err, ShouldBeNil)

				Convey("And I should not be able to find it", func() {
					var dbProject db.Project
					err := db.FindById(db.Projects(appDb.Session), result.Id, &dbProject)

					So(err.Error(), ShouldEqual, "not found")
				})
			})
		})
	})
}
