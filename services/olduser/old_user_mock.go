package olduser

import (
	"github.com/coaraujo/go-processor/domains"
	"github.com/stretchr/testify/mock"
)

//OldUserMock is a mock for User
type OldUserMock struct {
	mock.Mock
}

//Initialize is a mock for Initialize
func (o *OldUserMock) Initialize() error {
	GetInstance()
	instance = o
	return nil
}

//Get is a mock for Get
func (o *OldUserMock) Get(id string) (*domains.User, error) {
	args := o.Called(id)
	return args.Get(0).(*domains.User), args.Error(1)
}

//Insert is a mock for Insert
func (o *OldUserMock) Insert(user *domains.User) (string, error) {
	args := o.Called(user)
	return args.String(0), args.Error(1)
}
