package db

import (
	"github.com/pstuart2/go-teamcity"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type BuildType struct {
	Id           string   `bson:"_id" json:"id"`
	Name         string   `bson:"name" json:"name"`
	Description  string   `bson:"description" json:"description"`
	ProjectID    string   `bson:"projectId" json:"projectId"`
	Branches     []Branch `bson:"branches" json:"branches"`
	DashboardIds []string `bson:"dashboardIds" json:"dashboardIds"`
}

type Branch struct {
	Name   string  `bson:"name" json:"name"`
	Builds []Build `bson:"builds" json:"builds"`
}

type Build struct {
	Id         int                  `json:"id"`
	Number     string               `json:"number"`
	Status     teamcity.BuildStatus `json:"status"`
	StatusText string               `json:"statusText"`
	Progress   int                  `json:"progress"`
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

func (appDb *AppDb) UpdateBuildTypeBuilds(buildTypeId string, branches []Branch) (*BuildType, error) {
	now := appDb.now()

	change := mgo.Change{
		Update: bson.M{
			"$set": bson.M{
				"modifiedAt": now,
				"branches":   branches,
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

func (appDb *AppDb) AddDashboardToBuildTypes(buildTypeIds []string, dashboardId string) error {
	now := appDb.now()

	selector := bson.M{"_id": bson.M{"$in": buildTypeIds}}
	update := bson.M{"$set": bson.M{"modifiedAt": now}, "$push": bson.M{"dashboardIds": dashboardId}}

	_, err := BuildTypes(appDb.Session).UpdateAll(selector, update)
	return err
}

func (appDb *AppDb) RemoveDashboardFromBuildTypes(dashboardId string) error {
	now := appDb.now()

	selector := bson.M{"dashboardIds": dashboardId}
	update := bson.M{"$set": bson.M{"modifiedAt": now}, "$pull": bson.M{"dashboardIds": dashboardId}}

	_, err := BuildTypes(appDb.Session).UpdateAll(selector, update)
	return err
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

func (appDb *AppDb) DashboardBuildTypeList(dashboardId string) ([]BuildType, error) {
	var buildTypeList []BuildType

	if err := BuildTypes(appDb.Session).
		Find(bson.M{"deleted": bson.M{"$exists": false}, "dashboardIds": dashboardId}).
		All(&buildTypeList); err != nil {
		return nil, err
	}

	return buildTypeList, nil
}
