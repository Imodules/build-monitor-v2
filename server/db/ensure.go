package db

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

func Ensure(session *mgo.Session, log *logrus.Entry) error {
	if err := ensureUserCollection(Users(session), log); err != nil {
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
