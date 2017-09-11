package tc

import (
	"build-monitor-v2/server/db"

	"github.com/kapitanov/go-teamcity"
)

func (c *Server) RefreshProjects() error {
	projects, err := c.Tc.GetProjects()
	if err != nil {
		c.Log.Errorf("Failed to get projects from Team city: %v", err)
		return err
	}

	// TODO: Get project map, remove them as we process them

	c.Log.Infof("List of projects:\n")
	for _, project := range projects {
		if project.ID != "_Root" {
			dbProject := ProjectToDb(project)

			_, dbErr := c.Db.UpsertProject(dbProject)
			if dbErr != nil {
				c.Log.Errorf("Failed to upsert project. Id: %s, Name: %s", dbProject.Id, dbProject.Name)
			}
		}
	}

	// TODO: Delete from the db any projects left in the map

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
