package oldmetauser

import (
	"github.com/coaraujo/go-processor/domains"
	"github.com/stretchr/testify/mock"
)

//OldMetaUserMock is a mock for MetaUser
type OldMetaUserMock struct {
	mock.Mock
}

//Initialize is a mock for Initialize
func (g *OldMetaUserMock) Initialize() error {
	GetInstance()
	instance = g
	return nil
}

//Get is a mock for Get
func (o *OldMetaUserMock) Get(id string) (*domains.MetaUser, error) {
	args := o.Called(id)
	return args.Get(0).(*domains.MetaUser), args.Error(1)
}

//Insert is a mock for Insert
func (o *OldMetaUserMock) Insert(metauser *domains.MetaUser) (string, error) {
	args := o.Called(metauser)
	return args.String(0), args.Error(1)
}
