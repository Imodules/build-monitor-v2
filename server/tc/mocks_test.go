package tc_test

import (
	"build-monitor-v2/server/db"
	"github.com/kapitanov/go-teamcity"
	"github.com/stretchr/testify/mock"
)

type ITcClientMock struct {
	mock.Mock
}

func (m *ITcClientMock) GetProjects() ([]teamcity.Project, error) {
	args := m.Called()

	return args.Get(0).([]teamcity.Project), args.Error(1)
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
