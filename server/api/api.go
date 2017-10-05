package api

import (
	"build-monitor-v2/server/cfg"
	"build-monitor-v2/server/db"

	"fmt"

	"net/http"

	"context"
	"time"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

type IServer interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Use(middleware ...echo.MiddlewareFunc)

	Static(prefix, root string) *echo.Route
	Routes() []*echo.Route
	Group(prefix string, m ...echo.MiddlewareFunc) (g *echo.Group)
	Start(address string) error
	Shutdown(ctx context.Context) error

	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

type IGroup interface {
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

type IAppDb interface {
	CreateUser(username, email, password string) (*db.User, error)
	FindUserByLogin(usernameOrEmail string, password string) (*db.User, error)
	FindUserById(id string) (*db.User, error)
	LogUserLogin(user *db.User)

	ProjectList() ([]db.Project, error)
	BuildTypeList() ([]db.BuildType, error)

	DashboardList() ([]db.Dashboard, error)
	FindDashboardById(id string) (*db.Dashboard, error)
	UpsertDashboard(dashboard db.Dashboard) (*db.Dashboard, error)
	DeleteDashboard(id string) error

	AddDashboardToBuildTypes(buildTypeIds []string, dashboardId string) error
	RemoveDashboardFromBuildTypes(dashboardId string) error
	DashboardBuildTypeList(dashboardId string) ([]db.BuildType, error)
}

type ITcServer interface {
	Refresh()
}

type Server struct {
	Config    *cfg.Config
	Log       *logrus.Entry
	DbSession *mgo.Session
	Server    IServer
	TcServer  ITcServer
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func Create(log *logrus.Entry, config *cfg.Config, session *mgo.Session, tc ITcServer) *Server {
	return &Server{
		Config:    config,
		Log:       log,
		DbSession: session,
		Server:    echo.New(),
		TcServer:  tc,
	}
}

func (s *Server) Start() error {
	return s.Server.Start(fmt.Sprintf(":%d", s.Config.Port))
}

func (s *Server) Shutdown() {
	s.Log.Info("Stopping")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Server.Shutdown(ctx); err != nil {
		s.Log.Panicf("Failed to shutdown api server: %v", err)
	}

	s.Log.Info("Stopped")
}

func (s *Server) Refresh(ctx echo.Context) error {
	s.Log.Info("Refreshing")

	s.TcServer.Refresh()

	return nil
}
