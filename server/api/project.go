package api

import (
	"net/http"

	"github.com/labstack/echo"
)

func (s *Server) Projects(ctx echo.Context) error {
	appDb := getAppDb(ctx)

	projects, err := appDb.ProjectList()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, projects)
}
