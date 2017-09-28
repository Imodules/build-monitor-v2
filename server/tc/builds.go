package tc

import (
	"build-monitor-v2/server/db"

	"github.com/pstuart2/go-teamcity"
)

var GetRunningBuilds = func(c *Server) error {
	// httpAuth/app/rest/builds?locator=running:true
	return nil
}

var GetBuildHistory = func(c *Server) error {
	dashboards, dlErr := c.Db.DashboardList()
	if dlErr != nil {
		return dlErr
	}

	var btIdsList []string
	for _, d := range dashboards {
		for _, b := range d.BuildConfigs {
			if !contains(btIdsList, b.Id) {
				btIdsList = append(btIdsList, b.Id)
			}
		}
	}

	for _, buildTypeId := range btIdsList {
		builds, err := c.Tc.GetBuildsForBuildType(buildTypeId, 1000)
		if err != nil {
			c.Log.Errorf("Failed to get builds for buildType: %s, Error: %v", buildTypeId, err)
			continue
		}

		branchMap := make(map[string]*db.Branch)
		for _, build := range builds {
			var branch *db.Branch

			if val, ok := branchMap[build.BranchName]; ok {
				branch = val
			} else {
				branchMap[build.BranchName] = &db.Branch{Name: build.BranchName, Builds: []db.Build{}}
				branch = branchMap[build.BranchName]
			}

			branch.Builds = append(branch.Builds, BuildToDb(build))
		}

		if len(branchMap) > 0 {
			_, updateErr := c.Db.UpdateBuildTypeBuilds(buildTypeId, branchMapToArray(branchMap))
			if updateErr != nil {
				c.Log.Errorf("Failed to update db builds for buildType: %s, Error: %v", buildTypeId, updateErr)
			}
		}
	}

	return nil
}

func contains(s []string, k string) bool {
	for _, a := range s {
		if a == k {
			return true
		}
	}
	return false
}

func BuildToDb(p teamcity.Build) db.Build {
	return db.Build{
		Id:         p.ID,
		Number:     p.Number,
		Status:     p.Status,
		StatusText: p.StatusText,
		Progress:   p.Progress,
	}
}

func branchMapToArray(branches map[string]*db.Branch) []db.Branch {
	var arr []db.Branch
	for _, branch := range branches {
		arr = append(arr, *branch)
	}

	return arr
}
