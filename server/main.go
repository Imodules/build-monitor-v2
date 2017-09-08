package main

import (
	"build-monitor-v2/server/cfg"

	"build-monitor-v2/server/api"

	"build-monitor-v2/server/db"

	"build-monitor-v2/server/tc"

	"github.com/ian-kent/gofigure"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"os"
	"os/signal"
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

	server := api.Create(
		log.WithField("component", "api"),
		&config,
		session,
	)

	if err := server.Setup(); err != nil {
		log.Fatalf("Failed to setup server: %v", err)
	}

	tcMonitor := tc.NewServer(log.WithField("component", "tcMonitor"), &config)
	if err := tcMonitor.Start(); err != nil {
		log.Fatalf("Failed to start Teamcity monitor: %v", err)
	}

	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	waitForShutdownSignal(log)

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
