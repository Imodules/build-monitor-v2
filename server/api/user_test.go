package api_test

import (
	"encoding/json"
	"testing"

	"build-monitor-v2/server/api"

	"build-monitor-v2/server/db"

	"net/http"

	"errors"

	"build-monitor-v2/server/cfg"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2/bson"
)

func TestUserSignUp(t *testing.T) {
	Convey("Given a server", t, func() {
		config := cfg.Config{JwtSecret: "this world"}
		s := api.Server{Config: &config}

		Convey("With an invalid sign up request", func() {
			c, rec := createTestPostRequest("/signup", []byte("invalid"))

			resultErr := s.SignUp(c)
			So(resultErr, ShouldBeNil)

			Convey("It should return http.StatusBadRequest", func() {
				So(rec.Code, ShouldEqual, http.StatusBadRequest)

				Convey("And have an ErrorResponse", func() {
					var resp api.ErrorResponse
					err := json.Unmarshal(rec.Body.Bytes(), &resp)

					So(err, ShouldBeNil)
					So(resp.Message, ShouldNotBeEmpty)
				})
			})
		})

		Convey("With a valid sign up request", func() {
			r := &api.SignUpRequest{
				Username: "new username to sign up",
				Email:    "this is the email for it",
				Password: "not that creative",
			}

			requestJson, _ := json.Marshal(r)

			c, rec := createTestPostRequest("/signup", requestJson)

			mockDb := new(IAppDbMock)
			c.Set(dbKey, mockDb)

			Convey("When the user is successfully created", func() {
				dbUser := db.User{
					DbObject: db.DbObject{Id: bson.NewObjectId()},
					Username: r.Username,
					Email:    "the user email",
					Password: "a hashed password",
				}

				mockDb.On("CreateUser", r.Username, r.Email, r.Password).Return(&dbUser, nil)

				resultErr := s.SignUp(c)
				So(resultErr, ShouldBeNil)

				Convey("It should insert into the database and log the user login", func() {
					mockDb.AssertExpectations(t)

					Convey("And return http.StatusCreated", func() {
						So(rec.Code, ShouldEqual, http.StatusCreated)

						var resultUser db.User
						err := json.Unmarshal(rec.Body.Bytes(), &resultUser)
						So(err, ShouldBeNil)

						So(resultUser.Id.Hex(), ShouldEqual, dbUser.Id.Hex())
						So(resultUser.Username, ShouldEqual, dbUser.Username)
						So(resultUser.Email, ShouldEqual, dbUser.Email)

						Convey("And return a valid auth token", func() {
							So(resultUser.Token, ShouldNotBeEmpty)

							claims, err := getClaims(resultUser.Token, config.JwtSecret)

							So(err, ShouldBeNil)
							So(claims, ShouldNotBeNil)
							So(claims.Valid(), ShouldBeNil)

							Convey("And not return a password", func() {
								So(resultUser.Password, ShouldBeEmpty)
							})
						})
					})
				})

			})

			Convey("When the user is a duplicate", func() {
				mockDb.On("CreateUser", r.Username, r.Email, r.Password).Return(nil, db.DuplicateUser)

				resultErr := s.SignUp(c)
				So(resultErr, ShouldBeNil)

				Convey("It should attempt user creation", func() {
					mockDb.AssertExpectations(t)

					Convey("And return http.StatusConflict", func() {
						So(rec.Code, ShouldEqual, http.StatusConflict)

						var resp api.ErrorResponse
						err := json.Unmarshal(rec.Body.Bytes(), &resp)

						So(err, ShouldBeNil)
						So(resp.Message, ShouldEqual, db.DuplicateUser.Error())
					})
				})
			})

			Convey("When the user is missing a field", func() {
				mockDb.On("CreateUser", r.Username, r.Email, r.Password).Return(nil, db.MissingUserField)

				resultErr := s.SignUp(c)
				So(resultErr, ShouldBeNil)

				Convey("It should attempt user creation", func() {
					mockDb.AssertExpectations(t)

					Convey("And return http.StatusBadRequest", func() {
						So(rec.Code, ShouldEqual, http.StatusBadRequest)

						var resp api.ErrorResponse
						err := json.Unmarshal(rec.Body.Bytes(), &resp)

						So(err, ShouldBeNil)
						So(resp.Message, ShouldEqual, db.MissingUserField.Error())
					})
				})
			})

			Convey("When there is a general error creating the user", func() {
				mockDb.On("CreateUser", r.Username, r.Email, r.Password).Return(nil, errors.New("Wat?"))

				resultErr := s.SignUp(c)
				So(resultErr, ShouldBeNil)

				Convey("It should attempt user creation", func() {
					mockDb.AssertExpectations(t)

					Convey("And return http.StatusInternalServerError", func() {
						So(rec.Code, ShouldEqual, http.StatusInternalServerError)

						var resp api.ErrorResponse
						err := json.Unmarshal(rec.Body.Bytes(), &resp)

						So(err, ShouldBeNil)
						So(resp.Message, ShouldEqual, "Failed to create user")
					})
				})
			})
		})

	})
}

func TestUserLogin(t *testing.T) {
	Convey("Given a server", t, func() {
		config := cfg.Config{JwtSecret: "this world"}
		s := api.Server{Config: &config}

		Convey("With an invalid login request", func() {
			c, rec := createTestPostRequest("/login", []byte("invalid"))

			resultErr := s.Login(c)
			So(resultErr, ShouldBeNil)

			Convey("It should return http.StatusBadRequest", func() {
				So(rec.Code, ShouldEqual, http.StatusBadRequest)

				Convey("And have an ErrorResponse", func() {
					var resp api.ErrorResponse
					err := json.Unmarshal(rec.Body.Bytes(), &resp)

					So(err, ShouldBeNil)
					So(resp.Message, ShouldNotBeEmpty)
				})
			})
		})

		Convey("With a valid login request", func() {
			r := &api.LoginRequest{
				Username: "new username to sign up",
				Password: "not that creative",
			}

			loginRequest, _ := json.Marshal(r)

			c, rec := createTestPostRequest("/login", loginRequest)

			mockDb := new(IAppDbMock)
			c.Set(dbKey, mockDb)

			Convey("When the user is successfully logged in", func() {
				dbUser := db.User{
					DbObject: db.DbObject{Id: bson.NewObjectId()},
					Username: r.Username,
					Email:    "the user email",
					Password: "a hashed password",
				}

				mockDb.On("FindUserByLogin", r.Username, r.Password).Return(&dbUser, nil)
				mockDb.On("LogUserLogin", &dbUser).Once()

				resultErr := s.Login(c)
				So(resultErr, ShouldBeNil)

				Convey("It should return http.StatusOK and log the user login", func() {
					mockDb.AssertExpectations(t)

					So(rec.Code, ShouldEqual, http.StatusOK)

					var resultUser db.User
					err := json.Unmarshal(rec.Body.Bytes(), &resultUser)
					So(err, ShouldBeNil)

					So(resultUser.Id.Hex(), ShouldEqual, dbUser.Id.Hex())
					So(resultUser.Username, ShouldEqual, dbUser.Username)
					So(resultUser.Email, ShouldEqual, dbUser.Email)

					Convey("And return a valid auth token", func() {
						So(resultUser.Token, ShouldNotBeEmpty)

						claims, err := getClaims(resultUser.Token, config.JwtSecret)

						So(err, ShouldBeNil)
						So(claims, ShouldNotBeNil)
						So(claims.Valid(), ShouldBeNil)

						Convey("And not return a password", func() {
							So(resultUser.Password, ShouldBeEmpty)
						})
					})
				})

			})

			Convey("When the user is fails to login", func() {
				mockDb.On("FindUserByLogin", r.Username, r.Password).Return(nil, errors.New("Not found"))

				resultErr := s.Login(c)
				So(resultErr, ShouldBeNil)

				Convey("It should return http.StatusUnauthorized", func() {
					So(rec.Code, ShouldEqual, http.StatusUnauthorized)

					var resp api.ErrorResponse
					err := json.Unmarshal(rec.Body.Bytes(), &resp)

					So(err, ShouldBeNil)
					So(resp.Message, ShouldEqual, "Invalid username / password combination")
				})

			})

		})

	})

}

func TestUserReAuthenticate(t *testing.T) {
	Convey("Given a user attempts to re-authenticate", t, func() {
		username := "TheUserName"
		email := "The user's email"

		dbUser := &db.User{
			DbObject: db.DbObject{Id: bson.NewObjectId()},
			Username: username,
			Email:    email,
			Password: "hashed password",
		}

		c, rec := createTestGetRequest("/api/authenticate")

		mockDb := new(IAppDbMock)
		c.Set(dbKey, mockDb)

		config := cfg.Config{JwtSecret: "this world"}
		s := api.Server{Config: &config}

		setClaims(c, dbUser)

		Convey("When the user is found", func() {
			mockDb.On("FindUserById", dbUser.Id.Hex()).Return(dbUser, nil)
			mockDb.On("LogUserLogin", dbUser).Once()

			Convey("And the generateToken succeeds", func() {
				resultErr := s.ReAuthenticate(c)
				So(resultErr, ShouldBeNil)

				Convey("It should return http.StatusOK and the user object", func() {
					mockDb.AssertExpectations(t)

					So(rec.Code, ShouldEqual, http.StatusOK)
					var resultUser db.User
					err := json.Unmarshal(rec.Body.Bytes(), &resultUser)
					So(err, ShouldBeNil)

					So(resultUser.Id.Hex(), ShouldEqual, dbUser.Id.Hex())
					So(resultUser.Username, ShouldEqual, dbUser.Username)
					So(resultUser.Email, ShouldEqual, dbUser.Email)

					Convey("With a new auth token", func() {
						So(resultUser.Token, ShouldNotBeEmpty)

						claims, err := getClaims(resultUser.Token, config.JwtSecret)

						So(err, ShouldBeNil)
						So(claims, ShouldNotBeNil)
						So(claims.Valid(), ShouldBeNil)

						Convey("And not return a password", func() {
							So(resultUser.Password, ShouldBeEmpty)
						})
					})
				})
			})
		})

		Convey("When the user is not found", func() {

			mockDb.On("FindUserById", dbUser.Id.Hex()).Return(nil, db.UserNotFound)

			resultErr := s.ReAuthenticate(c)
			So(resultErr, ShouldBeNil)

			Convey("It should return http.StatusNotFound and an error", func() {
				mockDb.AssertExpectations(t)

				So(rec.Code, ShouldEqual, http.StatusNotFound)
				errorResult, _ := json.Marshal(api.ErrorResponse{Message: db.UserNotFound.Error()})
				So(rec.Body.String(), ShouldEqual, string(errorResult))
			})

		})

		Convey("When there is an unknown error", func() {

			mockDb.On("FindUserById", dbUser.Id.Hex()).Return(nil, errors.New("Oh dang"))

			resultErr := s.ReAuthenticate(c)
			So(resultErr, ShouldBeNil)

			Convey("It should return http.StatusInternalServerError and an error", func() {
				mockDb.AssertExpectations(t)

				So(rec.Code, ShouldEqual, http.StatusInternalServerError)
				errorResult, _ := json.Marshal(api.ErrorResponse{Message: "There was an unknown error while trying to find the user"})
				So(rec.Body.String(), ShouldEqual, string(errorResult))
			})

		})

	})

}
