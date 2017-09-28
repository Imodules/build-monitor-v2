package tc

import (
	"build-monitor-v2/server/db"

	"github.com/pstuart2/go-teamcity"
)

var RefreshBuildTypes = func(c *Server) error {
	projectMap, projErr := projectMap(c.Db)
	if projErr != nil {
		c.Log.Errorf("Failed to get project list from database: %v", projErr)
		return projErr
	}

	if len(projectMap) == 0 {
		return nil
	}

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
		if _, ok := projectMap[buildType.ProjectID]; ok {
			dbBuildType := BuildTypeToDb(buildType)

			_, dbErr := c.Db.UpsertBuildType(dbBuildType)
			if dbErr != nil {
				c.Log.Errorf("Failed to upsert buildType. Id: %s, Name: %s", dbBuildType.Id, dbBuildType.Name)
			}

			delete(dbBuildTypeMap, dbBuildType.Id)
		}
	}

	for _, buildType := range dbBuildTypeMap {
		c.Db.DeleteBuildType(buildType.Id)
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
