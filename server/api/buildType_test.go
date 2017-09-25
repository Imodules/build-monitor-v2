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

func TestServer_BuildTypes(t *testing.T) {
	Convey("Given a server", t, func() {
		config := cfg.Config{JwtSecret: "this world"}
		s := api.Server{Config: &config}

		c, rec := createTestGetRequest("/api/buildTypes")

		mockDb := new(IAppDbMock)
		c.Set(dbKey, mockDb)

		Convey("When there are buildTypes", func() {
			buildTypes := []db.BuildType{
				{Id: "Build_type_01", Name: "Build Type 01", ProjectID: "Project_01"},
				{Id: "Build_type_02", Name: "Build Type 02", ProjectID: "Project_02"},
				{Id: "Build_type_03", Name: "Build Type 03", ProjectID: "Project_02"},
				{Id: "Build_type_04", Name: "Build Type 04", ProjectID: "Project_01"},
			}

			mockDb.On("BuildTypeList").Return(buildTypes, nil)

			resultErr := s.BuildTypes(c)
			So(resultErr, ShouldBeNil)

			Convey("It should query the database", func() {
				mockDb.AssertExpectations(t)

				Convey("And return http.StatusOK", func() {
					So(rec.Code, ShouldEqual, http.StatusOK)

					var result []db.BuildType
					err := json.Unmarshal(rec.Body.Bytes(), &result)
					So(err, ShouldBeNil)

					So(len(result), ShouldEqual, 4)
					So(result[0].Id, ShouldEqual, "Build_type_01")
					So(result[1].Id, ShouldEqual, "Build_type_02")
					So(result[2].Id, ShouldEqual, "Build_type_03")
					So(result[3].Id, ShouldEqual, "Build_type_04")

					So(result[0].ProjectID, ShouldEqual, "Project_01")
					So(result[1].ProjectID, ShouldEqual, "Project_02")
					So(result[2].ProjectID, ShouldEqual, "Project_02")
					So(result[3].ProjectID, ShouldEqual, "Project_01")
				})
			})
		})

		Convey("When the database errors", func() {
			expectedError := errors.New("this is some bad mojo")
			mockDb.On("BuildTypeList").Return(nil, expectedError)

			resultErr := s.BuildTypes(c)
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
