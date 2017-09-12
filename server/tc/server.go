package tc

import (
	"build-monitor-v2/server/cfg"

	"time"

	"build-monitor-v2/server/db"

	"github.com/kapitanov/go-teamcity"
	"github.com/sirupsen/logrus"
)

type ITcClient interface {
	GetProjects() ([]teamcity.Project, error)
}

type IDb interface {
	UpsertProject(r db.Project) (*db.Project, error)
	ProjectList() ([]db.Project, error)
	DeleteProject(id string) error
}

type Server struct {
	Tc                       ITcClient
	Db                       IDb
	Log                      *logrus.Entry
	ProjectPollInterval      time.Duration
	BuildPollInterval        time.Duration
	RunningBuildPollInterval time.Duration
	stop                     chan bool
	stopped                  chan bool
}

func NewServer(log *logrus.Entry, c *cfg.Config, appDb IDb) Server {
	return Server{
		Tc:                       teamcity.NewClient(c.TcUrl, teamcity.GuestAuth()),
		Db:                       appDb,
		Log:                      log,
		ProjectPollInterval:      getIntervalDuration(log, "TcProjectPollInterval", c.TcProjectPollInterval),
		BuildPollInterval:        getIntervalDuration(log, "TcBuildPollInterval", c.TcBuildPollInterval),
		RunningBuildPollInterval: getIntervalDuration(log, "TcRunningBuildPollInterval", c.TcRunningBuildPollInterval),
	}
}

func (c *Server) Start() error {
	// Refresh projects on start to ensure we are able to connect and read from the server
	if err := RefreshProjects(c); err != nil {
		return err
	}

	c.stop = make(chan bool)
	c.stopped = make(chan bool)

	// Now start our monitor
	go monitor(c)

	return nil
}

func (c *Server) Shutdown() {
	c.stop <- true

	select {
	case <-c.stopped:
		c.Log.Info("Stopped")
		break
	case <-time.After(time.Second * 10):
		c.Log.Error("Failed to stop after 10 seconds")
	}
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

	for shouldStop == false {
		select {
		//case m := <-msgs:
		//	err := ci.c.process(&m, db, ci.influx, Log.WithField("deliveryTag", m.DeliveryTag))
		//	ci.handleProcessResult(ch, &m, err)
		case shouldStop = <-c.stop:
			c.Log.Info("Stopping")
			break
		case <-time.After(c.ProjectPollInterval):
			RefreshProjects(c)
		}
	}

	c.stopped <- true
}
