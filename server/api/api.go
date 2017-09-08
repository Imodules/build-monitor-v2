package api

import (
	"build-monitor-v2/server/cfg"
	"build-monitor-v2/server/db"

	"fmt"

	"net/http"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

type IServer interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Use(middleware ...echo.MiddlewareFunc)

	Static(prefix, root string)
	Routes() []*echo.Route
	Group(prefix string, m ...echo.MiddlewareFunc) (g *echo.Group)
	Start(address string) error

	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
}

type IGroup interface {
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc)
}

type IAppDb interface {
	CreateUser(username, email, password string) (*db.User, error)
	FindUserByLogin(usernameOrEmail string, password string) (*db.User, error)
	FindUserById(id string) (*db.User, error)
	LogUserLogin(user *db.User)
}

type Server struct {
	Config    *cfg.Config
	Log       *logrus.Entry
	DbSession *mgo.Session
	Server    IServer
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func Create(log *logrus.Entry, config *cfg.Config, session *mgo.Session) *Server {
	return &Server{
		Config:    config,
		Log:       log,
		DbSession: session,
		Server:    echo.New(),
	}
}

func (s *Server) Start() error {
	return s.Server.Start(fmt.Sprintf(":%d", s.Config.Port))
}
