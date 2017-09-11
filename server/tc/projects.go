package tc

import (
	"build-monitor-v2/server/db"

	"github.com/kapitanov/go-teamcity"
)

var RefreshProjects = func(c *Server) error {
	projects, err := c.Tc.GetProjects()
	if err != nil {
		c.Log.Errorf("Failed to get projects from Team city: %v", err)
		return err
	}

	dbProjectMap, pmErr := projectMap(c.Db)
	if pmErr != nil {
		c.Log.Errorf("Failed to get projects from Team city: %v", err)
		return pmErr
	}

	c.Log.Infof("List of projects:")
	for _, project := range projects {
		if project.ID != "_Root" {
			dbProject := ProjectToDb(project)

			_, dbErr := c.Db.UpsertProject(dbProject)
			if dbErr != nil {
				c.Log.Errorf("Failed to upsert project. Id: %s, Name: %s", dbProject.Id, dbProject.Name)
			}

			delete(dbProjectMap, dbProject.Id)
		}
	}

	for _, project := range dbProjectMap {
		c.Db.DeleteProject(project.Id)
	}

	return nil
}

func projectMap(appDb IDb) (map[string]db.Project, error) {
	projects, err := appDb.ProjectList()
	if err != nil {
		return nil, err
	}

	projectMap := make(map[string]db.Project)
	for _, v := range projects {
		projectMap[v.Id] = v
	}

	return projectMap, nil
}

func ProjectToDb(p teamcity.Project) db.Project {
	return db.Project{
		Id:              p.ID,
		Name:            p.Name,
		Description:     p.Description,
		ParentProjectID: p.ParentProjectID,
	}
}
