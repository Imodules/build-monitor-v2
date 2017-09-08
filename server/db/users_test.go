package db_test

import (
	"testing"

	"build-monitor-v2/server/cfg"
	"build-monitor-v2/server/db"

	"time"

	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2/bson"
)

func TestCreateUser(t *testing.T) {

	Convey("Given an appDb instance with new user information", t, func() {
		c := cfg.Config{PasswordSalt: "something here"}
		log := logrus.WithField("test", "TestCreateUser")

		appDb := db.Create(dbSession, &c, log)

		Convey("When validateUser succeeds", func() {

			now := time.Now()

			username := "My Cool User 01"
			email := "My Cool Email 01"
			password := "Some C00l Pass 01"

			Convey("And username and email do not already exist", func() {
				createdUser, err := appDb.CreateUser(username, email, password)

				Convey("Then a new user record is returned without an error", func() {

					So(err, ShouldBeNil)

					So(createdUser, ShouldNotBeNil)
					So(createdUser.Id.Hex(), ShouldNotEqual, "")
					So(createdUser.CreatedAt.Unix(), ShouldAlmostEqual, now.Unix(), 1.0)
					So(createdUser.ModifiedAt.Unix(), ShouldEqual, createdUser.CreatedAt.Unix())
					So(createdUser.LastLoginAt.Unix(), ShouldEqual, createdUser.CreatedAt.Unix())

					Convey("And it was inserted into the appDb", func() {
						var queriedUser db.User
						userErr := db.Users(dbSession).Find(bson.M{"_id": createdUser.Id}).One(&queriedUser)

						So(userErr, ShouldBeNil)
						So(queriedUser.Id.Hex(), ShouldEqual, createdUser.Id.Hex())
					})

				})

			})

			Convey("And username is a duplicate", func() {
				createdUser, err := appDb.CreateUser(username, "different", "more diff")

				Convey("Then DuplicateUser error is returned", func() {
					So(err, ShouldEqual, db.DuplicateUser)

					Convey("And a user record is not created", func() {
						So(createdUser, ShouldBeNil)
					})
				})

			})

			Convey("And email is a duplicate", func() {
				createdUser, err := appDb.CreateUser("different", email, "more diff")

				Convey("Then DuplicateUser error is returned", func() {
					So(err, ShouldEqual, db.DuplicateUser)

					Convey("And a user record is not created", func() {
						So(createdUser, ShouldBeNil)
					})
				})

			})
		})

		Convey("When username is missing", func() {
			createdUser, err := appDb.CreateUser("", "different", "more diff")

			Convey("Then MissingUserField error is returned", func() {
				So(err, ShouldEqual, db.MissingUserField)

				Convey("And a user record is not created", func() {
					So(createdUser, ShouldBeNil)
				})
			})

		})

		Convey("When email is missing", func() {
			createdUser, err := appDb.CreateUser("this is a good one", "", "more diff")

			Convey("Then MissingUserField error is returned", func() {
				So(err, ShouldEqual, db.MissingUserField)

				Convey("And a user record is not created", func() {
					So(createdUser, ShouldBeNil)
				})
			})

		})

		Convey("When password is missing", func() {
			createdUser, err := appDb.CreateUser("this is a good one", "oy!", "")

			Convey("Then MissingUserField error is returned", func() {
				So(err, ShouldEqual, db.MissingUserField)

				Convey("And a user record is not created", func() {
					So(createdUser, ShouldBeNil)
				})
			})

		})

	})

}

func TestUsers_FindUserByLogin(t *testing.T) {

	Convey("When searching for a user for logging in", t, func() {
		c := cfg.Config{PasswordSalt: "something here"}
		log := logrus.WithField("test", "TestUsers_FindUserByLogin")

		appDb := db.Create(dbSession, &c, log)

		Convey("When a username and password match a user", func() {
			user, userErr := appDb.CreateUser("psLogint1", "pslogint1@fwe.com", "cool1pass")
			So(userErr, ShouldBeNil)

			foundUser, err := appDb.FindUserByLogin("psLogint1", "cool1pass")

			Convey("It should not return an error", func() {
				So(err, ShouldBeNil)

				Convey("It should return that user object", func() {
					So(foundUser, ShouldNotBeNil)
					So(foundUser.Id, ShouldEqual, user.Id)
				})
			})
		})

		Convey("When an email and password match a user", func() {
			user, userErr := appDb.CreateUser("psLogint2", "pslogint2@fwe.com", "cool2pass")
			So(userErr, ShouldBeNil)

			foundUser, err := appDb.FindUserByLogin("pslogint2@fwe.com", "cool2pass")

			Convey("It not return an error", func() {
				So(err, ShouldBeNil)

				Convey("It should return that user object", func() {
					So(foundUser, ShouldNotBeNil)
					So(foundUser.Id, ShouldEqual, user.Id)
				})

			})

		})

		Convey("When a username or email match but the password does not", func() {
			_, errInsert := appDb.CreateUser("psLogint3", "pslogint3@fwe.com", "cool3pass")
			So(errInsert, ShouldBeNil)

			_, err := appDb.FindUserByLogin("psLogint3", "cool0pass")

			Convey("It should return a UserNotFound error", func() {
				So(err, ShouldEqual, db.UserNotFound)
			})

		})

		Convey("When a username or email is not found", func() {
			_, insertErr := appDb.CreateUser("psLogint4", "pslogint4@fwe.com", "cool4pass")

			So(insertErr, ShouldBeNil)

			_, err := appDb.FindUserByLogin("psLogint4no", "cool4pass")

			Convey("It should return a UserNotFound error", func() {
				So(err, ShouldEqual, db.UserNotFound)
			})

		})

	})

}

func TestFindUserById(t *testing.T) {
	Convey("Given an appDb", t, func() {
		c := cfg.Config{PasswordSalt: "something here"}
		log := logrus.WithField("test", "TestFindUserById")

		appDb := db.Create(dbSession, &c, log)

		Convey("When a user exists with that Id", func() {
			user, userErr := appDb.CreateUser("psFindId1", "psFindId1@fwe.com", "cool1pass")
			So(userErr, ShouldBeNil)

			foundUser, err := appDb.FindUserById(user.Id.Hex())

			Convey("It should not return an error", func() {
				So(err, ShouldBeNil)

				Convey("It should return that user object", func() {
					So(foundUser, ShouldNotBeNil)
					So(foundUser.Id, ShouldEqual, user.Id)
				})
			})
		})

		Convey("When a user with that id does not exist", func() {
			user, err := appDb.FindUserById(bson.NewObjectId().Hex())

			Convey("It should return a UserNotFound error", func() {
				So(user, ShouldBeNil)
				So(err, ShouldEqual, db.UserNotFound)
			})

		})
	})
}

func TestUsers_LogUserLogin(t *testing.T) {
	Convey("Given an appDb", t, func() {
		c := cfg.Config{}
		log := logrus.WithField("test", "TestUsers_LogUserLogin")

		appDb := db.Create(dbSession, &c, log)

		Convey("When LogUserLogin is called with a valid user", func() {
			user, userErr := appDb.CreateUser("psLastLogin", "psLastLogin@fwe.com", "cool2pass")
			So(userErr, ShouldBeNil)

			now := time.Now()

			appDb.LogUserLogin(user)

			Convey("It should update the LastLoginAt to the current time", func() {
				var queriedUser db.User
				userErr2 := db.Users(dbSession).Find(bson.M{"_id": user.Id}).One(&queriedUser)
				So(userErr2, ShouldBeNil)

				So(queriedUser.LastLoginAt.Unix(), ShouldAlmostEqual, now.Unix(), 1.0)
			})
		})

		Convey("When LogUserLogin is called with an invalid user", func() {
			user := db.User{DbObject: db.DbObject{Id: bson.NewObjectId()}}
			appDb.LogUserLogin(&user)

			Convey("It should not create the user", func() {
				var queriedUser db.User
				userErr2 := db.Users(dbSession).Find(bson.M{"_id": user.Id}).One(&queriedUser)
				So(userErr2.Error(), ShouldEqual, "not found")
			})
		})
	})
}
