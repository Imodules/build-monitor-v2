package db

import (
	"errors"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2"
)

var badEnsureCallCount int

func badEnsure(failCall int) func(c *mgo.Collection) error {
	badEnsureCallCount = 0

	return func(c *mgo.Collection) error {
		if badEnsureCallCount == failCall {
			return errors.New(fmt.Sprintf("Nope: %d / %d!", failCall, badEnsureCallCount))
		}

		badEnsureCallCount++
		return nil
	}
}

func TestEnsure_Ensure(t *testing.T) {
	dbSession, _ := mgo.Dial("mongodb://localhost/build-monitor-v2-test")
	defer dbSession.Close()

	log := logrus.WithField("test", "ensure_internal_test.go")

	Convey("When ensureUsername fails", t, func() {
		origEnsure := ensureUsername
		ensureUsername = badEnsure(0)
		defer func() { ensureUsername = origEnsure }()

		Convey("It should successfully ensure the indexes on the database", func() {
			err := Ensure(dbSession, log)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "Nope: 0 / 0!")
		})

	})

	Convey("When ensureEmail fails", t, func() {
		origEnsure := ensureEmail
		ensureEmail = badEnsure(0)
		defer func() { ensureEmail = origEnsure }()

		Convey("It should successfully ensure the indexes on the database", func() {
			err := Ensure(dbSession, log)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "Nope: 0 / 0!")
		})

	})

	Convey("When ensureTeamCityId fails", t, func() {
		origEnsure := ensureTeamCityId
		ensureTeamCityId = badEnsure(0)
		defer func() { ensureTeamCityId = origEnsure }()

		Convey("It should successfully ensure the indexes on the database", func() {
			err := Ensure(dbSession, log)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "Nope: 0 / 0!")
		})

	})

	Convey("When ensureDeleted fails", t, func() {
		origEnsure := ensureDeleted
		ensureDeleted = badEnsure(0)
		defer func() { ensureDeleted = origEnsure }()

		Convey("It should successfully ensure the indexes on the database", func() {
			err := Ensure(dbSession, log)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "Nope: 0 / 0!")
		})

	})
}
