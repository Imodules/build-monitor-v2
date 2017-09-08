package db

import (
	"time"

	"build-monitor-v2/server/cfg"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type AppDb struct {
	Session      *mgo.Session
	PasswordSalt string
	Log          *logrus.Entry
	now          func() time.Time
}

type DbObject struct {
	Id         bson.ObjectId `bson:"_id,omitempty" json:"id"`
	CreatedAt  time.Time     `bson:"createdAt" json:"createdAt"`
	ModifiedAt time.Time     `bson:"modifiedAt" json:"modifiedAt"`
}

type Owner struct {
	Id       bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Username string        `json:"username"`
}

var getId = bson.NewObjectId

func Create(s *mgo.Session, c *cfg.Config, log *logrus.Entry, now func() time.Time) *AppDb {
	d := AppDb{
		Session:      s,
		PasswordSalt: c.PasswordSalt,
		Log:          log,
		now:          now,
	}

	return &d
}

func FindById(c *mgo.Collection, id string, o interface{}) error {
	err := c.Find(bson.M{"_id": id, "deleted": bson.M{"$exists": false}}).One(o)
	if err != nil {
		return err
	}

	return nil
}

func (appDb *AppDb) Delete(c *mgo.Collection, id string) error {
	return c.UpdateId(id, bson.M{"$set": bson.M{"deleted": appDb.now()}})
}

func (appDb *AppDb) setCreated(do *DbObject) {
	do.Id = getId()
	do.CreatedAt = appDb.now()
	do.ModifiedAt = do.CreatedAt
}
