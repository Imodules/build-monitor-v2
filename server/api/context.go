package api

import (
	"strings"
	"time"

	"build-monitor-v2/server/db"

	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
)

const (
	tokenKey  = "app.token"
	loggerKey = "app.logger"
	dbKey     = "app.Db"
)

func getLogger(ctx echo.Context) *logrus.Entry {
	obj := ctx.Get(loggerKey)
	if obj == nil {
		return logrus.NewEntry(logrus.StandardLogger())
	}

	return obj.(*logrus.Entry)
}

func getToken(ctx echo.Context) *jwt.Token {
	obj := ctx.Get(tokenKey)
	if obj == nil {
		return nil
	}

	return obj.(*jwt.Token)
}

func getClaims(ctx echo.Context) *JWTClaims {
	obj := getToken(ctx)
	if obj == nil {
		return nil
	}

	return obj.Claims.(*JWTClaims)
}

func getAppDb(ctx echo.Context) IAppDb {
	obj := ctx.Get(dbKey)
	if obj == nil {
		return nil
	}

	return obj.(IAppDb)
}

func isAssetRequest(r *http.Request) bool {
	return strings.Index(r.URL.Path, "/assets/") == 0
}

func isNotApiRequest(r *http.Request) bool {
	return strings.Index(r.URL.Path, "/api") != 0
}

func getSetupRequestHandler(s *Server) func(f echo.HandlerFunc) echo.HandlerFunc {
	return func(f echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			req := ctx.Request()

			if isAssetRequest(req) {
				return f(ctx)
			}

			if isNotApiRequest(req) {
				return ctx.File(s.Config.ClientPath + "/index.html")
			}

			requestId := uuid.NewRandom().String()

			logger := s.Log.WithFields(logrus.Fields{
				"method":     req.Method,
				"path":       req.URL.Path,
				"request_id": requestId,
			})
			ctx.Set(loggerKey, logger)

			session := s.DbSession.Copy()
			defer session.Close()

			Db := db.Create(session, s.Config, logger, time.Now)
			ctx.Set(dbKey, Db)

			startTime := time.Now()
			defer func() {
				rsp := ctx.Response()
				logger.WithFields(logrus.Fields{
					"status_code":  rsp.Status,
					"runtime_nano": time.Since(startTime).Nanoseconds(),
				}).Info("Finished request")
			}()

			logger.WithFields(logrus.Fields{
				"user_agent":     req.UserAgent(),
				"content_length": req.ContentLength,
			}).Info("Starting request")

			// we have to do this b/c if not the final error handler will not
			// in the chain of middleware. It will be called after meaning that the
			// response won't be set properly.
			err := f(ctx)
			if err != nil {
				ctx.Error(err)
			}
			return err
		}
	}
}
