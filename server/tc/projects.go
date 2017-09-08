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

	return nil
}

func ProjectToDb(p teamcity.Project) db.Project {
	return db.Project{
		Id:              p.ID,
		Name:            p.Name,
		Description:     p.Description,
		ParentProjectID: p.ParentProjectID,
	}
}
