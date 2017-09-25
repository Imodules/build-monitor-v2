package api

import (
	"net/http"

	"github.com/labstack/echo"
)

func (s *Server) BuildTypes(ctx echo.Context) error {
	appDb := getAppDb(ctx)

	buildTypes, err := appDb.BuildTypeList()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, buildTypes)
}
