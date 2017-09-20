package db

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

func Ensure(session *mgo.Session, log *logrus.Entry) error {
	if err := ensureUserCollection(Users(session), log); err != nil {
		return err
	}

	if err := ensureProjectCollection(Projects(session), log); err != nil {
		return err
	}

	if err := ensureBuildTypeCollection(BuildTypes(session), log); err != nil {
		return err
	}

	return nil
}

func ensureUserCollection(c *mgo.Collection, log *logrus.Entry) error {
	if err := ensureUsername(c); err != nil {
		log.Error("Failed calling ensureUsername: ", err)
		return err
	}

	if err := ensureEmail(c); err != nil {
		log.Error("Failed calling ensureEmail: ", err)
		return err
	}

	return nil
}

func ensureProjectCollection(c *mgo.Collection, log *logrus.Entry) error {
	if err := ensureDeleted(c); err != nil {
		log.Error("Failed calling ensureDeleted: ", err)
		return err
	}

	return nil
}

func ensureBuildTypeCollection(c *mgo.Collection, log *logrus.Entry) error {
	if err := ensureProjectId(c); err != nil {
		log.Error("Failed calling ensureProjectId: ", err)
		return err
	}

	return nil
}

var ensureUsername = func(c *mgo.Collection) error {
	index := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	return c.EnsureIndex(index)
}

var ensureEmail = func(c *mgo.Collection) error {
	index := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	return c.EnsureIndex(index)
}

var ensureDeleted = func(c *mgo.Collection) error {
	index := mgo.Index{
		Key:         []string{"deleted"},
		Unique:      false,
		DropDups:    false,
		Background:  true,
		Sparse:      true,
		ExpireAfter: 60 * 60 * 24 * 7, // 1 week
	}
	return c.EnsureIndex(index)
}

var ensureProjectId = func(c *mgo.Collection) error {
	index := mgo.Index{
		Key:        []string{"projectId"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	}
	return c.EnsureIndex(index)
}
