package tc_test

import (
	"build-monitor-v2/server/db"

	"github.com/pstuart2/go-teamcity"
	"github.com/stretchr/testify/mock"
)

type ITcClientMock struct {
	mock.Mock
}

func (m *ITcClientMock) GetProjects() ([]teamcity.Project, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]teamcity.Project), args.Error(1)
}

func (m *ITcClientMock) GetBuildTypes() ([]teamcity.BuildType, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]teamcity.BuildType), args.Error(1)
}

type IDbMock struct {
	mock.Mock
}

func (m *IDbMock) UpsertProject(r db.Project) (*db.Project, error) {
	args := m.Called(r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*db.Project), args.Error(1)
}

func (m *IDbMock) ProjectList() ([]db.Project, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]db.Project), args.Error(1)
}

func (m *IDbMock) DeleteProject(id string) error {
	args := m.Called(id)

	return args.Error(0)
}

func (m *IDbMock) UpsertBuildType(r db.BuildType) (*db.BuildType, error) {
	args := m.Called(r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*db.BuildType), args.Error(1)
}

func (m *IDbMock) BuildTypeList() ([]db.BuildType, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]db.BuildType), args.Error(1)
}

func (m *IDbMock) DeleteBuildType(id string) error {
	args := m.Called(id)

	return args.Error(0)
}
