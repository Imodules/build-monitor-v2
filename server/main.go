package main

import (
	"build-monitor-v2/server/cfg"

	"build-monitor-v2/server/api"

	"build-monitor-v2/server/db"

	"build-monitor-v2/server/tc"

	"os"
	"os/signal"
	"time"

	"github.com/ian-kent/gofigure"
	"github.com/jasonlvhit/gocron"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

func main() {

	log := logrus.WithField("component", "main")

	log.Info("Loading config...")

	config, cfgErr := cfg.Load(gofigure.Gofigure)
	if cfgErr != nil {
		log.Fatalf("Failed to load configs: %v", cfgErr)
	}

	session := setupDatabase(log, config)
	defer session.Close()

	tcLog := log.WithField("component", "tcMonitor")
	tcDb := db.Create(session.Copy(), &config, tcLog, time.Now)

	tcMonitor := tc.NewServer(tcLog, &config, tcDb)
	if err := tcMonitor.Start(); err != nil {
		log.Fatalf("Failed to start Teamcity monitor: %v", err)
	}

	gocron.Every(1).Hour().Do(func() { tcMonitor.Refresh() })

	server := api.Create(
		log.WithField("component", "api"),
		&config,
		session,
		&tcMonitor,
	)

	if err := server.Setup(); err != nil {
		log.Fatalf("Failed to setup server: %v", err)
	}

	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	cronChannel := gocron.Start()

	waitForShutdownSignal(log)

	cronChannel <- true

	server.Shutdown()
	tcMonitor.Shutdown()
}

func setupDatabase(log *logrus.Entry, config cfg.Config) *mgo.Session {
	log.Info("Setting up database.")
	appDbMasterSession, err := mgo.Dial(config.Db)
	if err != nil {
		log.Fatal("Failed to dial appDbMasterSession: " + err.Error())
	}

	appDbMasterSession.SetMode(mgo.Monotonic, true)

	db.Ensure(appDbMasterSession, log)

	return appDbMasterSession
}

func waitForShutdownSignal(log *logrus.Entry) {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)

	log.Print(">>> Running. To exit press CTRL+C")

	<-sigChannel
	signal.Stop(sigChannel)
	log.Println("> Shutdown!")
}
