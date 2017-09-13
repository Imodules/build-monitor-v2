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
)

func TestServer_Projects(t *testing.T) {
	Convey("Given a server", t, func() {
		config := cfg.Config{JwtSecret: "this world"}
		s := api.Server{Config: &config}

		c, rec := createTestGetRequest("/api/projects")

		mockDb := new(IAppDbMock)
		c.Set(dbKey, mockDb)

		Convey("When there are projects", func() {
			projects := []db.Project{
				{Id: "Project_01", Name: "Project 01", ParentProjectID: "_Root"},
				{Id: "Project_02", Name: "Project 02", ParentProjectID: "_Root"},
				{Id: "Project_03", Name: "Project 03", ParentProjectID: "Project_02"},
				{Id: "Project_04", Name: "Project 04", ParentProjectID: "Project_01"},
			}

			mockDb.On("ProjectList").Return(projects, nil)

			resultErr := s.Projects(c)
			So(resultErr, ShouldBeNil)

			Convey("It should query the database", func() {
				mockDb.AssertExpectations(t)

				Convey("And return http.StatusOK", func() {
					So(rec.Code, ShouldEqual, http.StatusOK)

					var result []db.Project
					err := json.Unmarshal(rec.Body.Bytes(), &result)
					So(err, ShouldBeNil)

					So(len(result), ShouldEqual, 4)
					So(result[0].Id, ShouldEqual, "Project_01")
					So(result[1].Id, ShouldEqual, "Project_02")
					So(result[2].Id, ShouldEqual, "Project_03")
					So(result[3].Id, ShouldEqual, "Project_04")

					So(result[1].ParentProjectID, ShouldEqual, "_Root")
					So(result[2].ParentProjectID, ShouldEqual, "Project_02")
					So(result[3].ParentProjectID, ShouldEqual, "Project_01")
				})
			})
		})

		Convey("When the database errors", func() {
			expectedError := errors.New("this is some bad mojo")
			mockDb.On("ProjectList").Return(nil, expectedError)

			resultErr := s.Projects(c)
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
