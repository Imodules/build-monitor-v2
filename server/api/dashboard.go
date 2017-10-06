package api

import (
	"net/http"

	"build-monitor-v2/server/db"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func (s *Server) Dashboards(ctx echo.Context) error {
	appDb := getAppDb(ctx)

	dashboards, err := appDb.DashboardList()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, dashboards)
}

type DashboardDetails struct {
	Id               string            `json:"id"`
	Name             string            `json:"name"`
	ColumnCount      int               `json:"columnCount"`
	SuccessIcon      string            `json:"successIcon"`
	FailedIcon       string            `json:"failedIcon"`
	RunningIcon      string            `json:"runningIcon"`
	LeftDateFormat   string            `json:"leftDateFormat"`
	CenterDateFormat string            `json:"centerDateFormat"`
	RightDateFormat  string            `json:"rightDateFormat"`
	Details          []BuildTypeDetail `json:"details"`
}

type BuildTypeDetail struct {
	Id           string      `json:"id"`
	Name         string      `json:"name"`
	Abbreviation string      `json:"abbreviation"`
	IsRunning    bool        `bson:"isRunning" json:"isRunning"`
	Branches     []db.Branch `json:"branches"`
}

func (s *Server) DashboardDetails(ctx echo.Context) error {
	log := getLogger(ctx)
	appDb := getAppDb(ctx)

	id := ctx.Param("id")

	dashboard, dashErr := appDb.FindDashboardById(id)
	if dashErr != nil {
		log.Error("Failed to get the dashboard from the database", dashErr)
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: dashErr.Error()})
	}

	buildTypes, btErr := appDb.DashboardBuildTypeList(id)
	if btErr != nil {
		log.Error("Failed to get the buildTypes from the database", btErr)
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: btErr.Error()})
	}

	details := DashboardDetails{
		Id:               dashboard.Id,
		Name:             dashboard.Name,
		ColumnCount:      dashboard.ColumnCount,
		SuccessIcon:      dashboard.SuccessIcon,
		FailedIcon:       dashboard.FailedIcon,
		RunningIcon:      dashboard.RunningIcon,
		LeftDateFormat:   dashboard.LeftDateFormat,
		CenterDateFormat: dashboard.CenterDateFormat,
		RightDateFormat:  dashboard.RightDateFormat,
	}

	for _, c := range dashboard.BuildConfigs {
		buildType := findBuildType(c.Id, buildTypes)
		detail := BuildTypeDetail{Id: c.Id, Abbreviation: c.Abbreviation}

		if buildType != nil {
			detail.Name = buildType.Name
			detail.Branches = buildType.Branches
			detail.IsRunning = buildType.IsRunning
		}

		details.Details = append(details.Details, detail)
	}

	return ctx.JSON(http.StatusOK, details)
}

type UpdateDashboardRequest struct {
	Name             string           `json:"name"`
	ColumnCount      int              `json:"columnCount"`
	SuccessIcon      string           `json:"successIcon"`
	FailedIcon       string           `json:"failedIcon"`
	RunningIcon      string           `json:"runningIcon"`
	LeftDateFormat   string           `json:"leftDateFormat"`
	CenterDateFormat string           `json:"centerDateFormat"`
	RightDateFormat  string           `json:"rightDateFormat"`
	BuildConfigs     []db.BuildConfig `json:"buildConfigs"`
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
		Id:               bson.NewObjectId().Hex(),
		Name:             r.Name,
		ColumnCount:      r.ColumnCount,
		SuccessIcon:      r.SuccessIcon,
		FailedIcon:       r.FailedIcon,
		RunningIcon:      r.RunningIcon,
		LeftDateFormat:   r.LeftDateFormat,
		CenterDateFormat: r.CenterDateFormat,
		RightDateFormat:  r.RightDateFormat,
		Owner:            db.Owner{Id: bson.ObjectIdHex(claims.UserId), Username: claims.Username},
		BuildConfigs:     r.BuildConfigs,
	}

	dbDashboard, err := appDb.UpsertDashboard(dashboard)
	if err != nil {
		log.Error("Failed to insert dashboard into database", err)
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	addDashboardToBuildTypes(appDb, dbDashboard)

	return ctx.JSON(http.StatusCreated, dbDashboard)
}

func (s *Server) DeleteDashboard(ctx echo.Context) error {
	claims := getClaims(ctx)
	appDb := getAppDb(ctx)
	log := getLogger(ctx)

	id := ctx.Param("id")

	dashboard, _ := appDb.FindDashboardById(id)
	if dashboard.Owner.Id.Hex() != claims.UserId {
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse{Message: "You are not the owner"})
	}

	if err := appDb.RemoveDashboardFromBuildTypes(id); err != nil {
		log.Error("Failed to delete dashboard from the build types", err)
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
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
	if dbCheck.Owner.Id.Hex() != claims.UserId {
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse{Message: "You are not the owner"})
	}

	if err := appDb.RemoveDashboardFromBuildTypes(id); err != nil {
		log.Error("Failed to delete dashboard from the build types", err)
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	dashboard := db.Dashboard{
		Id:               id,
		Name:             r.Name,
		ColumnCount:      r.ColumnCount,
		SuccessIcon:      r.SuccessIcon,
		FailedIcon:       r.FailedIcon,
		RunningIcon:      r.RunningIcon,
		LeftDateFormat:   r.LeftDateFormat,
		CenterDateFormat: r.CenterDateFormat,
		RightDateFormat:  r.RightDateFormat,
		Owner:            db.Owner{Id: bson.ObjectIdHex(claims.UserId), Username: claims.Username},
		BuildConfigs:     r.BuildConfigs,
	}

	dbDashboard, err := appDb.UpsertDashboard(dashboard)
	if err != nil {
		log.Error("Failed to insert dashboard into database", err)
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	addDashboardToBuildTypes(appDb, dbDashboard)

	return ctx.JSON(http.StatusOK, dbDashboard)
}

func addDashboardToBuildTypes(appDb IAppDb, dbDashboard *db.Dashboard) {
	var ids []string
	for _, config := range dbDashboard.BuildConfigs {
		ids = append(ids, config.Id)
	}
	appDb.AddDashboardToBuildTypes(ids, dbDashboard.Id)
}

func findBuildType(id string, buildTypes []db.BuildType) *db.BuildType {
	for _, bt := range buildTypes {
		if bt.Id == id {
			return &bt
		}
	}

	return nil
}
