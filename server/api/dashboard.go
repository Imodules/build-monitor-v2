package api

import (
	"net/http"

	"build-monitor-v2/server/db"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func (s *Server) Dashboards(ctx echo.Context) error {
	appDb := getAppDb(ctx)

	claims := getClaims(ctx)

	dashboards, err := appDb.DashboardList(claims.UserId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, dashboards)
}

type UpdateDashboardRequest struct {
	Name         string   `json:"name"`
	BuildTypeIds []string `json:"buildTypeIds"`
}

func (s *Server) CreateDashboard(ctx echo.Context) error {
	r := new(UpdateDashboardRequest)
	if err := ctx.Bind(r); err != nil {
		return ctx.JSON(http.StatusInternalServerError, nil)
	}

	log := getLogger(ctx)
	claims := getClaims(ctx)
	appDb := getAppDb(ctx)

	dashboard := db.Dashboard{
		Id:           bson.NewObjectId().Hex(),
		Name:         r.Name,
		OwnerId:      claims.UserId,
		BuildTypeIds: r.BuildTypeIds,
	}

	dbDashboard, err := appDb.UpsertDashboard(dashboard)
	if err != nil {
		log.Error("Failed to insert dashboard into database", err)
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusCreated, dbDashboard)
}

func (s *Server) DeleteDashboard(ctx echo.Context) error {
	claims := getClaims(ctx)
	appDb := getAppDb(ctx)
	log := getLogger(ctx)

	id := ctx.Param("id")

	dashboard, _ := appDb.FindDashboardById(id)
	if dashboard.OwnerId != claims.UserId {
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse{Message: "You are not the owner"})
	}

	if err := appDb.DeleteDashboard(id); err != nil {
		log.Error("Failed to delete dashboard from database", err)
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (s *Server) UpdateDashboard(ctx echo.Context) error {
	r := new(UpdateDashboardRequest)
	if err := ctx.Bind(r); err != nil {
		return ctx.JSON(http.StatusInternalServerError, nil)
	}

	log := getLogger(ctx)
	claims := getClaims(ctx)
	appDb := getAppDb(ctx)

	id := ctx.Param("id")

	dbCheck, _ := appDb.FindDashboardById(id)
	if dbCheck.OwnerId != claims.UserId {
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse{Message: "You are not the owner"})
	}

	dashboard := db.Dashboard{
		Id:           id,
		Name:         r.Name,
		OwnerId:      claims.UserId,
		BuildTypeIds: r.BuildTypeIds,
	}

	dbDashboard, err := appDb.UpsertDashboard(dashboard)
	if err != nil {
		log.Error("Failed to insert dashboard into database", err)
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, dbDashboard)
}
