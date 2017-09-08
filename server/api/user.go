package api

import (
	"net/http"

	"build-monitor-v2/server/db"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type SignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) SignUp(ctx echo.Context) error {
	r := new(SignUpRequest)
	if err := ctx.Bind(r); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
	}

	log := getLogger(ctx)
	appDb := getAppDb(ctx)

	user, err := appDb.CreateUser(r.Username, r.Email, r.Password)
	if err != nil {
		if err == db.DuplicateUser {
			return ctx.JSON(http.StatusConflict, ErrorResponse{Message: err.Error()})
		}

		if err == db.MissingUserField {
			return ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		}

		log.Error("Failed to insert user into database", err)
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to create user"})
	}

	signedToken, tokenErr := GenerateToken(user, s.Config.JwtSecret)
	if tokenErr != nil {
		log.Error("Failed to generate token", err)
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to generate token"})
	}

	user.Token = signedToken

	log.WithFields(logrus.Fields{
		"_id":      user.Id.Hex(),
		"username": user.Username,
		"email":    user.Email,
	}).Info("User created")

	return ctx.JSON(http.StatusCreated, user)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Server) Login(ctx echo.Context) error {
	r := new(LoginRequest)
	if err := ctx.Bind(r); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
	}

	log := getLogger(ctx)
	appDb := getAppDb(ctx)

	user, err := appDb.FindUserByLogin(r.Username, r.Password)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Invalid username / password combination"})
	}

	signedToken, tokenErr := GenerateToken(user, s.Config.JwtSecret)
	if tokenErr != nil {
		log.Error("Failed to generate token", err)
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to generate token"})
	}

	user.Token = signedToken

	appDb.LogUserLogin(user)

	log.WithFields(logrus.Fields{
		"_id":      user.Id.Hex(),
		"username": user.Username,
		"email":    user.Email,
	}).Info("User logged in")

	return ctx.JSON(http.StatusOK, user)
}

func (s *Server) ReAuthenticate(ctx echo.Context) error {
	log := getLogger(ctx)
	appDb := getAppDb(ctx)

	token := getToken(ctx)
	claims := token.Claims.(*JWTClaims)

	user, err := appDb.FindUserById(claims.UserId)
	if err != nil {
		if err == db.UserNotFound {
			return ctx.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		}

		log.Error("Failed to lookup the user in the database", err)
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: "There was an unknown error while trying to find the user"})
	}

	signedToken, tokenErr := GenerateToken(user, s.Config.JwtSecret)
	if tokenErr != nil {
		log.Error("Failed to generate token", err)
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to generate token"})
	}

	user.Token = signedToken

	appDb.LogUserLogin(user)

	log.WithFields(logrus.Fields{
		"_id":      user.Id.Hex(),
		"username": user.Username,
		"email":    user.Email,
	}).Info("User logged in")

	return ctx.JSON(http.StatusOK, user)
}
