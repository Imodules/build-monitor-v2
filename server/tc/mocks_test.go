package tc_test

import (
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
