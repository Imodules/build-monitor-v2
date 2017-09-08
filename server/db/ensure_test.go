package db_test

import (
	"testing"

	"build-monitor-v2/server/db"

	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEnsure_Ensure(t *testing.T) {

	log := logrus.WithField("test", "ensure_test.go")

	Convey("When all calls are successful", t, func() {

		Convey("It should successfully ensure the indexes on the database", func() {
			err := db.Ensure(dbSession, log)
			So(err, ShouldBeNil)
		})

	})

}
