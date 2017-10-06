package tc

import (
	"build-monitor-v2/server/db"
	"sort"
	"time"

	"github.com/pstuart2/go-teamcity"
)

// GetRunningBuilds Gets the running builds
var GetRunningBuilds = func(c *Server, lastBuilds []teamcity.Build) []teamcity.Build {
	runningBuilds, err := c.Tc.GetRunningBuilds()
	if err != nil {
		c.Log.Errorf("Failed to get running builds, Error: %v", err)
		return lastBuilds
	}

	usefulBuilds := []teamcity.Build{}

	c.Log.Infof("Running: %v", runningBuilds)
	for _, b := range runningBuilds {
		bt, btErr := c.Db.FindBuildTypeById(b.BuildTypeID)
		if btErr != nil {
			c.Log.Errorf("Failed to get build type for: %s, Error: %v", b.BuildTypeID, btErr)
			continue
		}

		if len(bt.DashboardIds) == 0 {
			c.Log.Infof("No dashboards are monitoring this build, ignore")
			continue
		}

		usefulBuilds = append(usefulBuilds, b)

		ProcessRunningBuild(c, b, bt)
	}

	if len(lastBuilds) > 0 {
		for _, lb := range lastBuilds {
			if !isBuildInList(lb.ID, usefulBuilds) {
				bt, btErr := c.Db.FindBuildTypeById(lb.BuildTypeID)
				if btErr != nil {
					c.Log.Errorf("Failed to get build type for: %s, Error: %v", lb.BuildTypeID, btErr)
					continue
				}

				build, err := c.Tc.GetBuildByID(lb.ID)
				if err != nil {
					c.Log.Errorf("Failed to get the updated build for id: %d", lb.ID)
					continue
				}

				if updErr := ProcessRunningBuild(c, build, bt); updErr != nil {
					c.Log.Errorf("Failed to update builds for buildType: %s, Error: %v", bt.Id, updErr)
				}
			}
		}
	}

	return usefulBuilds
}

// ProcessRunningBuild merges a running build into build type and updates the db
var ProcessRunningBuild = func(c *Server, b teamcity.Build, bt *db.BuildType) error {
	index := indexOfBranch(b.BranchName, bt.Branches)
	if index == -1 {
		bt.Branches = append(bt.Branches, db.Branch{Name: b.BranchName})
		index = len(bt.Branches) - 1
	}

	newBuild := BuildToDb(b)

	if len(bt.Branches[index].Builds) == 0 {
		bt.Branches[index].Builds = []db.Build{newBuild}
	} else if bt.Branches[index].Builds[0].Id == newBuild.Id {
		bt.Branches[index].Builds[0] = newBuild
	} else {
		bt.Branches[index].Builds = append([]db.Build{newBuild}, bt.Branches[index].Builds...)
	}

	bt.Branches[index].Builds = cleanBuilds(bt.Branches[index].Builds)
	bt.Branches[index].IsRunning = isBranchRunning(bt.Branches[index].Builds)

	_, updErr := c.Db.UpdateBuildTypeBuilds(bt.Id, bt.Branches)
	return updErr
}

type buildsById []db.Build

func (s buildsById) Len() int {
	return len(s)
}
func (s buildsById) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s buildsById) Less(i, j int) bool {
	return s[j].Id < s[i].Id
}

func isBranchRunning(builds []db.Build) bool {
	for _, b := range builds {
		if b.Status == teamcity.StatusRunning {
			return true
		}
	}

	return false
}

func cleanBuilds(builds []db.Build) []db.Build {
	sort.Sort(buildsById(builds))

	if len(builds) > 12 {
		return builds[:12]
	}

	return builds
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
			branches := branchMapToArray(branchMap)
			for i, branch := range branches {
				branches[i].Builds = cleanBuilds(branch.Builds)
			}

			_, updateErr := c.Db.UpdateBuildTypeBuilds(buildTypeId, branches)
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
		StartDate:  p.StartDate,
		FinishDate: getCleanFinishDate(p.StartDate, p.FinishDate),
	}
}

func getCleanFinishDate(start, finish time.Time) time.Time {
	if start.After(finish) {
		return time.Now()
	}

	return finish
}

func branchMapToArray(branches map[string]*db.Branch) []db.Branch {
	var arr []db.Branch
	for _, branch := range branches {
		arr = append(arr, *branch)
	}

	return arr
}

func indexOfBranch(name string, branches []db.Branch) int {
	for k, v := range branches {
		if v.Name == name {
			return k
		}
	}

	return -1
}

func isBuildInList(id int, builds []teamcity.Build) bool {
	for _, v := range builds {
		if v.ID == id {
			return true
		}
	}

	return false
}
