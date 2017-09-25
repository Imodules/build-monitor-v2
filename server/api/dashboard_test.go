package api_test

import (
	"testing"

	"build-monitor-v2/server/api"

	"build-monitor-v2/server/cfg"

	"build-monitor-v2/server/db"

	"encoding/json"
	"net/http"

	"errors"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2/bson"
)

func TestServer_Dashboards(t *testing.T) {
	Convey("Given a server", t, func() {
		config := cfg.Config{JwtSecret: "this world"}
		s := api.Server{Config: &config}

		c, rec := createTestGetRequest("/api/dashboards")

		mockDb := new(IAppDbMock)
		c.Set(dbKey, mockDb)

		dbUser := &db.User{
			DbObject: db.DbObject{Id: bson.NewObjectId()},
			Username: "pstuart",
			Email:    "test@eample.com",
			Password: "hashed password",
		}

		setClaims(c, dbUser)

		Convey("When there are dashboards", func() {
			dashboards := []db.Dashboard{
				{Id: "Dashboard_01", Name: "Dashboard 01", OwnerId: dbUser.Id.Hex(), BuildTypeIds: []string{"a1", "b1", "c1"}},
				{Id: "Dashboard_02", Name: "Dashboard 02", OwnerId: dbUser.Id.Hex(), BuildTypeIds: []string{"a2", "b2", "c2"}},
				{Id: "Dashboard_03", Name: "Dashboard 03", OwnerId: dbUser.Id.Hex(), BuildTypeIds: []string{"a3", "b3", "c3"}},
				{Id: "Dashboard_04", Name: "Dashboard 04", OwnerId: dbUser.Id.Hex(), BuildTypeIds: []string{"a4", "b4", "c4"}},
			}

			mockDb.On("DashboardList", dbUser.Id.Hex()).Return(dashboards, nil)

			resultErr := s.Dashboards(c)
			So(resultErr, ShouldBeNil)

			Convey("It should query the database", func() {
				mockDb.AssertExpectations(t)

				Convey("And return http.StatusOK", func() {
					So(rec.Code, ShouldEqual, http.StatusOK)

					var result []db.Dashboard
					err := json.Unmarshal(rec.Body.Bytes(), &result)
					So(err, ShouldBeNil)

					So(len(result), ShouldEqual, 4)
					So(result[0].Id, ShouldEqual, "Dashboard_01")
					So(result[1].Id, ShouldEqual, "Dashboard_02")
					So(result[2].Id, ShouldEqual, "Dashboard_03")
					So(result[3].Id, ShouldEqual, "Dashboard_04")

					So(len(result[0].BuildTypeIds), ShouldEqual, 3)
					So(len(result[1].BuildTypeIds), ShouldEqual, 3)
					So(len(result[2].BuildTypeIds), ShouldEqual, 3)
					So(len(result[3].BuildTypeIds), ShouldEqual, 3)
				})
			})
		})

		Convey("When the database errors", func() {
			expectedError := errors.New("this is some bad mojo")
			mockDb.On("DashboardList", dbUser.Id.Hex()).Return(nil, expectedError)

			resultErr := s.Dashboards(c)
			So(resultErr, ShouldBeNil)

			Convey("It should return http.StatusInternalServerError and an error", func() {
				mockDb.AssertExpectations(t)

				So(rec.Code, ShouldEqual, http.StatusInternalServerError)
				errorResult, _ := json.Marshal(api.ErrorResponse{Message: expectedError.Error()})
				So(rec.Body.String(), ShouldEqual, string(errorResult))
			})
		})
	})
}

func TestServer_CreateDashboard(t *testing.T) {
	Convey("Given a server", t, func() {
		config := cfg.Config{JwtSecret: "this world"}
		s := api.Server{Config: &config}

		Convey("With an invalid request", func() {
			c, rec := createTestPostRequest("/api/dashboards", []byte{})

			err := s.CreateDashboard(c)
			Convey("It should not return an error", func() {
				So(err, ShouldBeNil)

				Convey("And return http.StatusInternalServerError", func() {

					So(rec.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})
		})

		Convey("With a valid request", func() {
			request := api.UpdateDashboardRequest{
				Name:         "This is me new dashboard",
				BuildTypeIds: []string{"db1", "db2"},
			}

			requestJson, _ := json.Marshal(request)

			c, rec := createTestPostRequest("/api/dashboards", requestJson)

			mockDb := new(IAppDbMock)
			c.Set(dbKey, mockDb)

			dbUser := &db.User{
				DbObject: db.DbObject{Id: bson.NewObjectId()},
				Username: "pstuart",
				Email:    "test@eample.com",
				Password: "hashed password",
			}

			setClaims(c, dbUser)

			dbDashboard := db.Dashboard{
				Id:           "from db",
				OwnerId:      dbUser.Id.Hex(),
				Name:         request.Name,
				BuildTypeIds: request.BuildTypeIds,
			}

			Convey("When the create succeeds", func() {
				mockDb.On("UpsertDashboard", mock.AnythingOfType("db.Dashboard")).Return(&dbDashboard, nil)

				resultErr := s.CreateDashboard(c)

				Convey("It should upsert the dashboard with the owner and a new id", func() {

					So(resultErr, ShouldBeNil)

					mockDb.AssertExpectations(t)

					dashboardToDb := mockDb.Calls[0].Arguments[0].(db.Dashboard)

					So(dashboardToDb.Id, ShouldNotBeEmpty)
					So(dashboardToDb.OwnerId, ShouldEqual, dbUser.Id.Hex())
					So(dashboardToDb.Name, ShouldEqual, request.Name)
					So(dashboardToDb.BuildTypeIds, ShouldContain, request.BuildTypeIds[0])
					So(dashboardToDb.BuildTypeIds, ShouldContain, request.BuildTypeIds[1])

					Convey("And return http.StatusCreated", func() {

						So(rec.Code, ShouldEqual, http.StatusCreated)

						Convey("And the dashboard", func() {
							var resultDashboard db.Dashboard
							err := json.Unmarshal(rec.Body.Bytes(), &resultDashboard)
							So(err, ShouldBeNil)

							expectedString, _ := json.Marshal(dbDashboard)
							So(rec.Body.String(), ShouldEqual, string(expectedString))

						})
					})
				})
			})

			Convey("When the create fails", func() {
				expectedErr := errors.New("what now")
				mockDb.On("UpsertDashboard", mock.AnythingOfType("db.Dashboard")).Return(nil, expectedErr)

				resultErr := s.CreateDashboard(c)

				Convey("It should upsert the dashboard with the owner and a new id", func() {

					So(resultErr, ShouldBeNil)

					mockDb.AssertExpectations(t)

					dashboardToDb := mockDb.Calls[0].Arguments[0].(db.Dashboard)

					So(dashboardToDb.Id, ShouldNotBeEmpty)
					So(dashboardToDb.OwnerId, ShouldEqual, dbUser.Id.Hex())
					So(dashboardToDb.Name, ShouldEqual, request.Name)
					So(dashboardToDb.BuildTypeIds, ShouldContain, request.BuildTypeIds[0])
					So(dashboardToDb.BuildTypeIds, ShouldContain, request.BuildTypeIds[1])

					Convey("And return http.StatusInternalServerError", func() {
						So(rec.Code, ShouldEqual, http.StatusInternalServerError)

						var resp api.ErrorResponse
						err := json.Unmarshal(rec.Body.Bytes(), &resp)

						So(err, ShouldBeNil)
						So(resp.Message, ShouldEqual, expectedErr.Error())
					})
				})
			})
		})
	})
}

func TestServer_DeleteDashboard(t *testing.T) {
	Convey("Given a server", t, func() {
		config := cfg.Config{JwtSecret: "this world"}
		s := api.Server{Config: &config}

		id := "98uie"
		c, rec := createTestDeleteRequest("/api/dashboards/" + id)

		c.SetParamNames("id")
		c.SetParamValues(id)

		mockDb := new(IAppDbMock)
		c.Set(dbKey, mockDb)

		dbUser := &db.User{
			DbObject: db.DbObject{Id: bson.NewObjectId()},
			Username: "pstuart",
			Email:    "test@eample.com",
			Password: "hashed password",
		}

		setClaims(c, dbUser)

		Convey("When the user owns the dashboard", func() {
			dashboard := db.Dashboard{
				Id:      "some id here",
				OwnerId: dbUser.Id.Hex(),
			}

			mockDb.On("FindDashboardById", id).Return(&dashboard, nil)

			Convey("And it successfully deletes from the db", func() {
				mockDb.On("DeleteDashboard", id).Return(nil)

				err := s.DeleteDashboard(c)
				So(err, ShouldBeNil)

				Convey("It should delete the dashboard", func() {
					mockDb.AssertExpectations(t)

					Convey("And return http.StatusOK", func() {

						So(rec.Code, ShouldEqual, http.StatusOK)
					})
				})
			})

			Convey("And it fails to delete from the db", func() {
				expectedErr := errors.New("this was expected")
				mockDb.On("DeleteDashboard", id).Return(expectedErr)

				err := s.DeleteDashboard(c)
				So(err, ShouldBeNil)

				mockDb.AssertExpectations(t)

				Convey("It should return http.StatusInternalServerError", func() {
					So(rec.Code, ShouldEqual, http.StatusInternalServerError)

					var resp api.ErrorResponse
					err := json.Unmarshal(rec.Body.Bytes(), &resp)

					So(err, ShouldBeNil)
					So(resp.Message, ShouldEqual, expectedErr.Error())
				})
			})
		})

		Convey("When the user does not own the dashboard", func() {
			dashboard := db.Dashboard{
				Id:      "some id here",
				OwnerId: "not this user",
			}

			mockDb.On("FindDashboardById", id).Return(&dashboard, nil)

			err := s.DeleteDashboard(c)
			So(err, ShouldBeNil)

			Convey("It should not delete from db", func() {

				mockDb.AssertExpectations(t)

				Convey("And return http.StatusUnauthorized", func() {
					So(rec.Code, ShouldEqual, http.StatusUnauthorized)
				})
			})
		})
	})
}

func TestServer_UpdateDashboard(t *testing.T) {
	Convey("Given a server", t, func() {
		config := cfg.Config{JwtSecret: "this world"}
		s := api.Server{Config: &config}

		id := "9iolk"

		Convey("With an invalid request", func() {
			c, rec := createTestPutRequest("/api/dashboards/"+id, []byte{})

			err := s.CreateDashboard(c)
			Convey("It should not return an error", func() {
				So(err, ShouldBeNil)

				Convey("And return http.StatusInternalServerError", func() {

					So(rec.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})
		})

		Convey("With a valid request", func() {
			request := api.UpdateDashboardRequest{
				Name:         "This is me new dashboard",
				BuildTypeIds: []string{"db1", "db2"},
			}

			requestJson, _ := json.Marshal(request)

			c, rec := createTestPostRequest("/api/dashboards/"+id, requestJson)

			c.SetParamNames("id")
			c.SetParamValues(id)

			mockDb := new(IAppDbMock)
			c.Set(dbKey, mockDb)

			dbUser := &db.User{
				DbObject: db.DbObject{Id: bson.NewObjectId()},
				Username: "pstuart",
				Email:    "test@eample.com",
				Password: "hashed password",
			}

			setClaims(c, dbUser)

			Convey("When we are the owher", func() {
				dbDashboard := db.Dashboard{
					Id:           "from db",
					OwnerId:      dbUser.Id.Hex(),
					Name:         request.Name,
					BuildTypeIds: request.BuildTypeIds,
				}

				Convey("And the update succeeds", func() {
					mockDb.On("FindDashboardById", id).Return(&dbDashboard, nil)
					mockDb.On("UpsertDashboard", mock.AnythingOfType("db.Dashboard")).Return(&dbDashboard, nil)

					resultErr := s.UpdateDashboard(c)

					Convey("It should upsert the dashboard with the owner and a new id", func() {

						So(resultErr, ShouldBeNil)

						mockDb.AssertExpectations(t)

						dashboardToDb := mockDb.Calls[1].Arguments[0].(db.Dashboard)

						So(dashboardToDb.Id, ShouldEqual, id)
						So(dashboardToDb.OwnerId, ShouldEqual, dbUser.Id.Hex())
						So(dashboardToDb.Name, ShouldEqual, request.Name)
						So(dashboardToDb.BuildTypeIds, ShouldContain, request.BuildTypeIds[0])
						So(dashboardToDb.BuildTypeIds, ShouldContain, request.BuildTypeIds[1])

						Convey("And return http.StatusOK", func() {

							So(rec.Code, ShouldEqual, http.StatusOK)

							Convey("And the dashboard", func() {
								var resultDashboard db.Dashboard
								err := json.Unmarshal(rec.Body.Bytes(), &resultDashboard)
								So(err, ShouldBeNil)

								expectedString, _ := json.Marshal(dbDashboard)
								So(rec.Body.String(), ShouldEqual, string(expectedString))

							})
						})
					})
				})

				Convey("And the update fails", func() {
					expectedErr := errors.New("what now")
					mockDb.On("FindDashboardById", id).Return(&dbDashboard, nil)
					mockDb.On("UpsertDashboard", mock.AnythingOfType("db.Dashboard")).Return(nil, expectedErr)

					resultErr := s.UpdateDashboard(c)

					Convey("It should upsert the dashboard with the owner and a new id", func() {

						So(resultErr, ShouldBeNil)

						mockDb.AssertExpectations(t)

						dashboardToDb := mockDb.Calls[1].Arguments[0].(db.Dashboard)

						So(dashboardToDb.Id, ShouldNotBeEmpty)
						So(dashboardToDb.OwnerId, ShouldEqual, dbUser.Id.Hex())
						So(dashboardToDb.Name, ShouldEqual, request.Name)
						So(dashboardToDb.BuildTypeIds, ShouldContain, request.BuildTypeIds[0])
						So(dashboardToDb.BuildTypeIds, ShouldContain, request.BuildTypeIds[1])

						Convey("And return http.StatusInternalServerError", func() {
							So(rec.Code, ShouldEqual, http.StatusInternalServerError)

							var resp api.ErrorResponse
							err := json.Unmarshal(rec.Body.Bytes(), &resp)

							So(err, ShouldBeNil)
							So(resp.Message, ShouldEqual, expectedErr.Error())
						})
					})
				})
			})

			Convey("When we are not the owner", func() {
				dbDashboard := db.Dashboard{
					Id:           "from db",
					OwnerId:      "some one else",
					Name:         request.Name,
					BuildTypeIds: request.BuildTypeIds,
				}

				mockDb.On("FindDashboardById", id).Return(&dbDashboard, nil)

				resultErr := s.UpdateDashboard(c)
				So(resultErr, ShouldBeNil)

				Convey("It should not update the dashboard", func() {
					mockDb.AssertExpectations(t)

					Convey("And return http.StatusUnauthorized", func() {
						So(rec.Code, ShouldEqual, http.StatusUnauthorized)

						var resp api.ErrorResponse
						err := json.Unmarshal(rec.Body.Bytes(), &resp)

						So(err, ShouldBeNil)
						So(resp.Message, ShouldEqual, "You are not the owner")
					})
				})
			})
		})
	})
}
