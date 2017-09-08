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

var getNow = time.Now
var getId = bson.NewObjectId

func Create(s *mgo.Session, c *cfg.Config, log *logrus.Entry) *AppDb {
	d := AppDb{
		Session:      s,
		PasswordSalt: c.PasswordSalt,
		Log:          log,
	}

	return &d
}

func setCreated(do *DbObject) {
	do.Id = getId()
	do.CreatedAt = getNow()
	do.ModifiedAt = do.CreatedAt
}
