package api

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// This is tested with integration tests
func (s *Server) Setup() error {

	setupMiddleware(s)
	setupRoutes(s)

	logRoutes(s)

	return nil
}

func setupMiddleware(s *Server) {
	s.Server.Use(getSetupRequestHandler(s))
	s.Server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{s.Config.AllowedOrigin},
	}))
	s.Server.Use(middleware.Static(s.Config.ClientPath))
}

func setupRoutes(s *Server) {

	s.Server.GET("*", func(c echo.Context) error {
		return c.File(s.Config.ClientPath + "index.html")
	})

	s.Server.POST("/api/signup", s.SignUp)
	s.Server.POST("/api/login", s.Login)

	requireClaims := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: jwt.SigningMethodHS256.Name,
		ContextKey:    tokenKey,
		Claims:        &JWTClaims{},
		SigningKey:    []byte(s.Config.JwtSecret),
	})

	secureGroup := s.Server.Group("/api", requireClaims)
	secureGroup.GET("/authenticate", s.ReAuthenticate)

	secureGroup.GET("/projects", s.Projects)
	secureGroup.GET("/buildTypes", s.BuildTypes)

	secureGroup.GET("/dashboards", s.Dashboards)
	secureGroup.POST("/dashboards", s.CreateDashboard)
	secureGroup.PUT("/dashboards/:id", s.UpdateDashboard)
	secureGroup.DELETE("/dashboards/:id", s.DeleteDashboard)
}

func logRoutes(s *Server) {
	routes := s.Server.Routes()
	for i := 0; i < len(routes); i++ {
		s.Log.Info(routes[i].Method + ": " + routes[i].Path)
	}
}
