package metauser

import (
	"github.com/coaraujo/go-processor/domains"
	"github.com/stretchr/testify/mock"
)

//MetaUserMock is a mock for MetaUser
type MetaUserMock struct {
	mock.Mock
}

//Initialize is a mock for Initialize
func (g *MetaUserMock) Initialize() error {
	GetInstance()
	instance = g
	return nil
}

//Get is a mock for Get
func (m *MetaUserMock) Get(id string) (*domains.MetaUser, error) {
	args := m.Called(id)
	return args.Get(0).(*domains.MetaUser), args.Error(1)
}

//Insert is a mock for Insert
func (m *MetaUserMock) Insert(user *domains.User, clientID string) (string, error) {
	args := m.Called(user, clientID)
	return args.String(0), args.Error(1)
}

//Update is a mock for Update
func (m *MetaUserMock) Update(newUser *domains.User, oldUser *domains.MetaUser, clientID string) error {
	args := m.Called(newUser, oldUser, clientID)
	return args.Error(0)
}

//Delete is a mock for  Delete
func (m *MetaUserMock) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
