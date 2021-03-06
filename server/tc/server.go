package tc

import (
	"build-monitor-v2/server/cfg"

	"time"

	"build-monitor-v2/server/db"

	"github.com/pstuart2/go-teamcity"
	"github.com/sirupsen/logrus"
)

type ITcClient interface {
	GetProjects() ([]teamcity.Project, error)
	GetBuildTypes() ([]teamcity.BuildType, error)
	GetBuildsForBuildType(id string, count int) ([]teamcity.Build, error)
	GetRunningBuilds() ([]teamcity.Build, error)
	GetBuildByID(id int) (teamcity.Build, error)
}

type IDb interface {
	UpsertProject(r db.Project) (*db.Project, error)
	ProjectList() ([]db.Project, error)
	DeleteProject(id string) error

	UpsertBuildType(r db.BuildType) (*db.BuildType, error)
	UpdateBuildTypeBuilds(buildTypeId string, branches []db.Branch) (*db.BuildType, error)
	BuildTypeList() ([]db.BuildType, error)
	DeleteBuildType(id string) error
	DashboardList() ([]db.Dashboard, error)
	FindBuildTypeById(id string) (*db.BuildType, error)
}

type Server struct {
	Tc                         ITcClient
	Db                         IDb
	Log                        *logrus.Entry
	TcPollInterval             time.Duration
	TcRunningBuildPollInterval time.Duration
	commands                   chan string
	stopped                    chan bool
}

func NewServer(log *logrus.Entry, c *cfg.Config, appDb IDb) Server {
	return Server{
		Tc:                         teamcity.NewClient(c.TcUrl, teamcity.GuestAuth()),
		Db:                         appDb,
		Log:                        log,
		TcPollInterval:             getIntervalDuration(log, "TcPollInterval", c.TcPollInterval),
		TcRunningBuildPollInterval: getIntervalDuration(log, "TcBuildPollInterval", c.TcRunningBuildPollInterval),
	}
}

func (c *Server) Start() error {
	// Refresh projects on start to ensure we are able to connect and read from the server
	if err := refresh(c); err != nil {
		return err
	}

	c.commands = make(chan string)
	c.stopped = make(chan bool)

	// Now start our monitor
	go monitor(c)

	return nil
}

func (c *Server) Shutdown() {
	c.commands <- "shutdown"

	select {
	case <-c.stopped:
		c.Log.Info("Stopped")
		break
	case <-time.After(time.Second * 10):
		c.Log.Error("Failed to stop after 10 seconds")
	}
}

func (c *Server) Refresh() {
	c.commands <- "refresh"
}

func getIntervalDuration(log *logrus.Entry, name, interval string) time.Duration {
	d, ciError := time.ParseDuration(interval)
	if ciError != nil {
		log.Panicf("Failed to parse duration %s = %s", name, interval)
	}

	return d
}

func monitor(c *Server) {
	c.Log.Info("Starting Teamcity monitor")
	shouldStop := false

	currentPollInterval := c.TcPollInterval
	runningBuilds := []teamcity.Build{}

	for shouldStop == false {
		select {
		case command := <-c.commands:
			switch command {
			case "shutdown":
				c.Log.Info("Stopping")
				shouldStop = true
				break
			case "refresh":
				refresh(c)
			}

		case <-time.After(currentPollInterval):
			runningBuilds = GetRunningBuilds(c, runningBuilds)
			if len(runningBuilds) == 0 {
				currentPollInterval = c.TcPollInterval
			} else {
				currentPollInterval = c.TcRunningBuildPollInterval
			}
		}
	}

	c.stopped <- true
}

func refresh(c *Server) error {
	if err := RefreshProjects(c); err != nil {
		return err
	}

	if err := RefreshBuildTypes(c); err != nil {
		return err
	}

	return GetBuildHistory(c)
}
