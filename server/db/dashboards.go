package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Dashboard struct {
	Id               string        `bson:"_id" json:"id"`
	Name             string        `bson:"name" json:"name"`
	ColumnCount      int           `bson:"columnCount" json:"columnCount"`
	Owner            Owner         `bson:"owner" json:"owner"`
	SuccessIcon      string        `bson:"successIcon" json:"successIcon"`
	FailedIcon       string        `bson:"failedIcon" json:"failedIcon"`
	RunningIcon      string        `bson:"runningIcon" json:"runningIcon"`
	LeftDateFormat   string        `bson:"leftDateFormat" json:"leftDateFormat"`
	CenterDateFormat string        `bson:"centerDateFormat" json:"centerDateFormat"`
	RightDateFormat  string        `bson:"rightDateFormat" json:"rightDateFormat"`
	BuildConfigs     []BuildConfig `bson:"buildConfigs" json:"buildConfigs"`
}

type BuildConfig struct {
	Id           string `bson:"_id" json:"id"`
	Abbreviation string `bson:"abbreviation" json:"abbreviation"`
}

func Dashboards(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("dashboards")
}

func (appDb *AppDb) UpsertDashboard(r Dashboard) (*Dashboard, error) {
	now := appDb.now()

	change := mgo.Change{
		Update: bson.M{
			"$set": bson.M{
				"modifiedAt":       now,
				"name":             r.Name,
				"columnCount":      r.ColumnCount,
				"successIcon":      r.SuccessIcon,
				"failedIcon":       r.FailedIcon,
				"runningIcon":      r.RunningIcon,
				"leftDateFormat":   r.LeftDateFormat,
				"centerDateFormat": r.CenterDateFormat,
				"rightDateFormat":  r.RightDateFormat,
				"owner":            r.Owner,
				"buildConfigs":     r.BuildConfigs,
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

func (appDb *AppDb) DashboardList() ([]Dashboard, error) {
	dashboardList := []Dashboard{}

	if err := Dashboards(appDb.Session).
		Find(bson.M{"deleted": bson.M{"$exists": false}}).
		Sort("name").
		Select(bson.M{
			"_id":              1,
			"name":             1,
			"owner":            1,
			"columnCount":      1,
			"successIcon":      1,
			"failedIcon":       1,
			"runningIcon":      1,
			"leftDateFormat":   1,
			"centerDateFormat": 1,
			"rightDateFormat":  1,
			"buildConfigs":     1,
		}).All(&dashboardList); err != nil {
		return nil, err
	}

	return dashboardList, nil
}
