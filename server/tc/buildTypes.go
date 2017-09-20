package tc

import (
	"build-monitor-v2/server/db"

	"github.com/kapitanov/go-teamcity"
)

// TODO: This should only refresh build types for projects we have.
var RefreshBuildTypes = func(c *Server) error {
	buildTypes, err := c.Tc.GetBuildTypes()
	if err != nil {
		c.Log.Errorf("Failed to get buildTypes from Team city: %v", err)
		return err
	}

	dbBuildTypeMap, pmErr := buildTypeMap(c.Db)
	if pmErr != nil {
		c.Log.Errorf("Failed to get buildTypes from Team city: %v", err)
		return pmErr
	}

	c.Log.Infof("List of buildTypes:")
	for _, buildType := range buildTypes {
		if buildType.ID != "_Root" {
			dbBuildType := BuildTypeToDb(buildType)

			_, dbErr := c.Db.UpsertBuildType(dbBuildType)
			if dbErr != nil {
				c.Log.Errorf("Failed to upsert buildType. Id: %s, Name: %s", dbBuildType.Id, dbBuildType.Name)
			}

			delete(dbBuildTypeMap, dbBuildType.Id)
		}
	}

	for _, project := range dbBuildTypeMap {
		c.Db.DeleteBuildType(project.Id)
	}

	return nil
}

func buildTypeMap(appDb IDb) (map[string]db.BuildType, error) {
	buildTypes, err := appDb.BuildTypeList()
	if err != nil {
		return nil, err
	}

	projectMap := make(map[string]db.BuildType)
	for _, v := range buildTypes {
		projectMap[v.Id] = v
	}

	return projectMap, nil
}

func BuildTypeToDb(p teamcity.BuildType) db.BuildType {
	return db.BuildType{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		ProjectID:   p.ProjectID,
	}
}
