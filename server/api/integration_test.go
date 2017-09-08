package api

import (
	"os"
	"testing"

	"build-monitor-v2/server/cfg"
	"build-monitor-v2/server/db"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo"
	. "github.com/smartystreets/goconvey/convey"
)

var s *Server
var config *cfg.Config
var user db.User

func TestMain(m *testing.M) {
	config = &cfg.Config{
		Port:         0,
		PasswordSalt: "integration-salt",
		JwtSecret:    "integration-secrt",
	}

	session, _ := mgo.Dial("mongodb://localhost/build-monitor-v2-integration")
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	session.DB("").DropDatabase()

	testLogger := logrus.StandardLogger().WithField("integration_test", true)
	db.Ensure(session, testLogger)

	s = Create(testLogger, config, session)
	s.Setup()

	os.Exit(m.Run())
}

func Test_000_MissingEndpoint(t *testing.T) {
	Convey("When an end point does not exist", t, func() {
		code, body := request(t, "GET", "/missing", nil)
		So(code, ShouldEqual, http.StatusNotFound)
		So(body, ShouldEqual, `{"message":"Not Found"}`)
	})
}

func Test_001_SignUp(t *testing.T) {
	Convey("When a user signs up", t, func() {
		payload := &SignUpRequest{
			Username: "paul-int-tester",
			Email:    "pstuart-int-tester@fwe.com",
			Password: "something not to share",
		}

		code, body := request(t, "POST", "/api/signup", payload)
		So(code, ShouldEqual, http.StatusCreated)

		suUser := new(db.User)
		extractPayload(t, body, suUser)

		So(suUser.Id, ShouldNotBeNil)
		So(suUser.Username, ShouldEqual, payload.Username)
		So(suUser.Email, ShouldEqual, payload.Email)
		So(suUser.Password, ShouldBeEmpty)
		So(suUser.Token, ShouldNotBeEmpty)
	})
}

func Test_003_Login(t *testing.T) {
	Convey("When a user logs in", t, func() {
		payload := &LoginRequest{
			Username: "paul-int-tester",
			Password: "something not to share",
		}

		code, body := request(t, "POST", "/api/login", payload)
		So(code, ShouldEqual, http.StatusOK)

		user = db.User{}
		extractPayload(t, body, &user)

		So(user.Id, ShouldNotBeNil)
		So(user.Username, ShouldEqual, payload.Username)
		So(user.Email, ShouldEqual, "pstuart-int-tester@fwe.com")
		So(user.Password, ShouldBeEmpty)
		So(user.Token, ShouldNotBeEmpty)
	})
}

func Test_004_Authenticate(t *testing.T) {
	Convey("When a user authenticates", t, func() {
		code, body := request(t, "GET", "/api/authenticate", nil)
		So(code, ShouldEqual, http.StatusOK)

		extractPayload(t, body, &user)

		So(user.Id, ShouldNotBeNil)
		So(user.Username, ShouldEqual, "paul-int-tester")
		So(user.Email, ShouldEqual, "pstuart-int-tester@fwe.com")
		So(user.Password, ShouldBeEmpty)
		So(user.Token, ShouldNotBeEmpty)
	})
}

// ------------------------------------------------------------------------------------------------
// Helpers
// ------------------------------------------------------------------------------------------------

func extractPayload(t *testing.T, body string, out interface{}) {
	err := json.Unmarshal([]byte(body), out)
	So(err, ShouldBeNil)
}

func request(t *testing.T, method, path string, body interface{}) (int, string) {
	req, _ := http.NewRequest(method, path, nil)

	if body != nil {
		bs, err := json.Marshal(body)
		if err != nil {
			t.FailNow()
		}

		req, _ = http.NewRequest(method, path, bytes.NewBuffer(bs))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}

	if user.Token != "" {
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+user.Token)
	}

	rsp := httptest.NewRecorder()
	s.Server.ServeHTTP(rsp, req)
	return rsp.Code, rsp.Body.String()
}
