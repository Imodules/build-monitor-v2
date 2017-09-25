package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Dashboard struct {
	Id           string   `bson:"_id" json:"id"`
	Name         string   `bson:"name" json:"name"`
	OwnerId      string   `bson:"ownerId" json:"ownerId"`
	BuildTypeIds []string `bson:"buildTypeIds" json:"buildTypeIds"`
}

func Dashboards(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("dashboards")
}

func (appDb *AppDb) UpsertDashboard(r Dashboard) (*Dashboard, error) {
	now := appDb.now()

	change := mgo.Change{
		Update: bson.M{
			"$set": bson.M{
				"modifiedAt":   now,
				"name":         r.Name,
				"ownerId":      r.OwnerId,
				"buildTypeIds": r.BuildTypeIds,
			},
			"$unset":       bson.M{"deleted": ""},
			"$setOnInsert": bson.M{"createdAt": now},
		},
		Upsert:    true,
		ReturnNew: true,
	}

	var dashboard Dashboard
	_, err := Dashboards(appDb.Session).Find(bson.M{
		"_id": r.Id,
	}).Apply(change, &dashboard)

	if err != nil {
		return nil, err
	}

	return &dashboard, nil
}

func (appDb *AppDb) DeleteDashboard(id string) error {
	return appDb.Delete(Dashboards(appDb.Session), id)
}

func (appDb *AppDb) FindDashboardById(id string) (*Dashboard, error) {
	var dashboard Dashboard
	if err := FindById(Dashboards(appDb.Session), id, &dashboard); err != nil {
		return nil, err
	}

	return &dashboard, nil
}

func (appDb *AppDb) DashboardList(ownerId string) ([]Dashboard, error) {
	dashboardList := []Dashboard{}

	if err := Dashboards(appDb.Session).
		Find(bson.M{"deleted": bson.M{"$exists": false}, "ownerId": ownerId}).
		Sort("name").
		Select(bson.M{
			"_id":          1,
			"name":         1,
			"ownerId":      1,
			"buildTypeIds": 1,
		}).All(&dashboardList); err != nil {
		return nil, err
	}

	return dashboardList, nil
}
