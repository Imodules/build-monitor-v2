package api

import (
	jwt "github.com/dgrijalva/jwt-go"
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

	s.Server.Static("/assets", s.Config.ClientPath)

	openApi := s.Server.Group("/api")
	openApi.POST("/signup", s.SignUp)
	openApi.POST("/login", s.Login)
	openApi.GET("/projects", s.Projects)
	openApi.GET("/buildTypes", s.BuildTypes)
	openApi.GET("/dashboards", s.Dashboards)
	openApi.GET("/dashboards/:id", s.DashboardDetails)

	requireClaims := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: jwt.SigningMethodHS256.Name,
		ContextKey:    tokenKey,
		Claims:        &JWTClaims{},
		SigningKey:    []byte(s.Config.JwtSecret),
	})

	secureApi := s.Server.Group("/api", requireClaims)
	secureApi.GET("/authenticate", s.ReAuthenticate)
	secureApi.POST("/dashboards", s.CreateDashboard)
	secureApi.PUT("/dashboards/:id", s.UpdateDashboard)
	secureApi.DELETE("/dashboards/:id", s.DeleteDashboard)
	secureApi.POST("/refresh", s.Refresh)
}

func logRoutes(s *Server) {
	routes := s.Server.Routes()
	for i := 0; i < len(routes); i++ {
		s.Log.Info(routes[i].Method + ": " + routes[i].Path)
	}
}
