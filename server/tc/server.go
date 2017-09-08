package tc

import (
	"build-monitor-v2/server/cfg"

	"time"

	"github.com/kapitanov/go-teamcity"
	"github.com/sirupsen/logrus"
)

type ITcClient interface {
	GetProjects() ([]teamcity.Project, error)
}

type Server struct {
	Tc                       ITcClient
	Log                      *logrus.Entry
	ProjectPollInterval      time.Duration
	BuildPollInterval        time.Duration
	RunningBuildPollInterval time.Duration
	stop                     chan bool
	stopped                  chan bool
}

func NewServer(log *logrus.Entry, c *cfg.Config) Server {
	return Server{
		Tc:                       teamcity.NewClient(c.TcUrl, teamcity.GuestAuth()),
		Log:                      log,
		ProjectPollInterval:      getIntervalDuration(log, "TcProjectPollInterval", c.TcProjectPollInterval),
		BuildPollInterval:        getIntervalDuration(log, "TcBuildPollInterval", c.TcBuildPollInterval),
		RunningBuildPollInterval: getIntervalDuration(log, "TcRunningBuildPollInterval", c.TcRunningBuildPollInterval),
	}
}

func (c *Server) Start() error {
	// Refresh projects on start to ensure we are able to connect and read from the server
	if err := c.refreshProjects(); err != nil {
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
			c.refreshProjects()
		}
	}

	c.stopped <- true
}

func (c *Server) refreshProjects() error {
	projects, err := c.Tc.GetProjects()
	if err != nil {
		c.Log.Errorf("Failed to get projects from Team city: %v", err)
		return err
	}

	c.Log.Infof("List of projects:\n")
	for _, project := range projects {
		c.Log.Infof(" * %s", project.ID)
	}

	return nil
}
