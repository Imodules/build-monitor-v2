package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type BuildType struct {
	Id          string `bson:"_id" json:"id"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	ProjectID   string `bson:"projectId" json:"projectId"`
	Paused      bool   `bson:"paused" json:"paused"`
}

func BuildTypes(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("buildTypes")
}

func (appDb *AppDb) UpsertBuildType(r BuildType) (*BuildType, error) {
	now := appDb.now()

	change := mgo.Change{
		Update: bson.M{
			"$set": bson.M{
				"modifiedAt":  now,
				"name":        r.Name,
				"description": r.Description,
				"projectId":   r.ProjectID,
				"paused":      r.Paused,
			},
			"$unset":       bson.M{"deleted": ""},
			"$setOnInsert": bson.M{"createdAt": now},
		},
		Upsert:    true,
		ReturnNew: true,
	}

	var buildType BuildType
	_, err := BuildTypes(appDb.Session).Find(bson.M{
		"_id": r.Id,
	}).Apply(change, &buildType)

	if err != nil {
		return nil, err
	}

	return &buildType, nil
}

func (appDb *AppDb) DeleteBuildType(id string) error {
	return appDb.Delete(BuildTypes(appDb.Session), id)
}

func (appDb *AppDb) BuildTypeList() ([]BuildType, error) {
	var buildTypeList []BuildType

	if err := BuildTypes(appDb.Session).
		Find(bson.M{"deleted": bson.M{"$exists": false}}).
		Sort("name").
		Select(bson.M{
			"_id":         1,
			"name":        1,
			"description": 1,
			"projectId":   1,
			"paused":      1,
		}).All(&buildTypeList); err != nil {
		return nil, err
	}

	return buildTypeList, nil
}
