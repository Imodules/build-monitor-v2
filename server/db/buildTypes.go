package db

import (
	"github.com/pstuart2/go-teamcity"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type BuildType struct {
	Id          string  `bson:"_id" json:"id"`
	Name        string  `bson:"name" json:"name"`
	Description string  `bson:"description" json:"description"`
	ProjectID   string  `bson:"projectId" json:"projectId"`
	Builds      []Build `bson:"builds" json:"builds"`
}

type Build struct {
	Id         int                  `json:"id"`
	Number     string               `json:"number"`
	Status     teamcity.BuildStatus `json:"status"`
	StatusText string               `json:"statusText"`
	Progress   int                  `json:"progress"`
	BranchName string               `json:"branchName"`
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

func (appDb *AppDb) UpdateBuilds(buildTypeId string, builds []Build) (*BuildType, error) {
	now := appDb.now()

	change := mgo.Change{
		Update: bson.M{
			"$set": bson.M{
				"modifiedAt": now,
				"builds":     builds,
			},
			"$unset": bson.M{"deleted": ""},
		},
		Upsert:    false,
		ReturnNew: true,
	}

	var buildType BuildType
	_, err := BuildTypes(appDb.Session).Find(bson.M{
		"_id": buildTypeId,
	}).Apply(change, &buildType)

	if err != nil {
		return nil, err
	}

	return &buildType, nil
}

func (appDb *AppDb) FindBuildTypeById(id string) (*BuildType, error) {
	var buildType BuildType
	if err := FindById(BuildTypes(appDb.Session), id, &buildType); err != nil {
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
