package api_test

import (
	"fmt"
	"time"

	"build-monitor-v2/server/api"
	"build-monitor-v2/server/db"

	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/mock"
)

const (
	tokenKey  = "app.token"
	loggerKey = "app.logger"
	dbKey     = "app.Db"
)

//region Help Methods
func getClaims(token string, secret string) (*api.JWTClaims, error) {
	result, err := jwt.ParseWithClaims(token, &api.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	return result.Claims.(*api.JWTClaims), nil
}

func setClaims(ctx echo.Context, user *db.User) {
	claims := &api.JWTClaims{
		UserId:   user.Id.Hex(),
		Username: user.Username,
		Email:    user.Email,
	}
	claims.ExpiresAt = time.Now().Add(time.Minute * 60).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ctx.Set(tokenKey, token)
}

func getToken(user *db.User) *jwt.Token {
	claims := &api.JWTClaims{
		UserId:   user.Id.Hex(),
		Username: user.Username,
		Email:    user.Email,
	}
	claims.ExpiresAt = time.Now().Add(time.Minute * 60).Unix()

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}

func createTestGetRequest(path string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req, _ := http.NewRequest(echo.GET, path, strings.NewReader(""))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return c, rec
}

func createTestPostRequest(path string, body []byte) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req, _ := http.NewRequest(echo.POST, path, strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return c, rec
}

//endregion

//region IServerMock Mock
type IServerMock struct {
	mock.Mock
}

func (m *IServerMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func (m *IServerMock) Use(middleware ...echo.MiddlewareFunc) {
	m.Called(middleware)
}

func (m *IServerMock) Static(prefix, root string) {
	m.Called(prefix, root)
}

func (em *IServerMock) Routes() []*echo.Route {
	args := em.Called()
	return args.Get(0).([]*echo.Route)
}

func (em *IServerMock) Group(prefix string, m ...echo.MiddlewareFunc) (g *echo.Group) {
	args := em.Called(prefix, m)
	return args.Get(0).(*echo.Group)
}

func (em *IServerMock) Start(address string) error {
	args := em.Called(address)
	return args.Error(0)
}

func (em *IServerMock) GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	em.Called(path, h, m)
}

func (em *IServerMock) POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	em.Called(path, h, m)
}

func (em *IServerMock) PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	em.Called(path, h, m)
}

func (em *IServerMock) PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	em.Called(path, h, m)
}

//endregion

//region IAppDbMock
type IAppDbMock struct {
	mock.Mock
}

func (m *IAppDbMock) CreateUser(username, email, password string) (*db.User, error) {
	args := m.Called(username, email, password)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*db.User), args.Error(1)
}

func (m *IAppDbMock) FindUserByLogin(usernameOrEmail string, password string) (*db.User, error) {
	args := m.Called(usernameOrEmail, password)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*db.User), args.Error(1)
}

func (m *IAppDbMock) FindUserById(id string) (*db.User, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*db.User), args.Error(1)
}

func (m *IAppDbMock) LogUserLogin(user *db.User) {
	m.Called(user)
}

//endregion
