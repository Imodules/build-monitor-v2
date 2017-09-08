package main

import (
	"build-monitor-v2/server/cfg"

	"build-monitor-v2/server/api"

	"build-monitor-v2/server/db"

	"github.com/ian-kent/gofigure"
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

	server := api.Create(
		log.WithField("component", "api"),
		&config,
		session,
	)

	if err := server.Setup(); err != nil {
		log.Fatalf("Failed to setup server: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
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
