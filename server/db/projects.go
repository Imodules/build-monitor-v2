package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Project struct {
	Id              string `bson:"_id" json:"id"`
	Name            string `bson:"name" json:"name"`
	Description     string `bson:"description" json:"description"`
	ParentProjectID string `bson:"parentProjectId" json:"parentProjectId"`
}

func Projects(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("projects")
}

func (appDb *AppDb) UpsertProject(r Project) (*Project, error) {
	now := appDb.now()

	change := mgo.Change{
		Update: bson.M{
			"$set": bson.M{
				"modifiedAt":      now,
				"name":            r.Name,
				"description":     r.Description,
				"parentProjectId": r.ParentProjectID,
			},
			"$unset":       bson.M{"deleted": ""},
			"$setOnInsert": bson.M{"createdAt": now},
		},
		Upsert:    true,
		ReturnNew: true,
	}

	var project Project
	_, err := Projects(appDb.Session).Find(bson.M{
		"_id": r.Id,
	}).Apply(change, &project)

	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (appDb *AppDb) DeleteProject(id string) error {
	return appDb.Delete(Projects(appDb.Session), id)
}

func (appDb *AppDb) ProjectList() ([]Project, error) {
	var projectList []Project

	if err := Projects(appDb.Session).
		Find(bson.M{"deleted": bson.M{"$exists": false}}).
		Sort("name").
		Select(bson.M{
			"_id":             1,
			"name":            1,
			"description":     1,
			"parentProjectId": 1,
		}).All(&projectList); err != nil {
		return nil, err
	}

	return projectList, nil
}
