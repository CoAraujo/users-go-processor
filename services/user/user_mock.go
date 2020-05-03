package user

import (
	"github.com/coaraujo/users-go-processor/domains"
	"github.com/stretchr/testify/mock"
)

//UserMock is a mock for User
type UserMock struct {
	mock.Mock
}

//Initialize is a mock for Initialize
func (u *UserMock) Initialize() error {
	GetInstance()
	instance = u
	return nil
}

//Get is a mock for Get
func (u *UserMock) Get(id string) (*domains.User, error) {
	args := u.Called(id)
	return args.Get(0).(*domains.User), args.Error(1)
}

//Insert is a mock for Insert
func (u *UserMock) Insert(user *domains.User) (string, error) {
	args := u.Called(user)
	return args.String(0), args.Error(1)
}

//Update is a mock for Update
func (u *UserMock) Update(newUser *domains.User, oldUser *domains.User) error {
	args := u.Called(newUser, oldUser)
	return args.Error(0)
}

//Delete is a mock for  Delete
func (u *UserMock) Delete(id string) error {
	args := u.Called(id)
	return args.Error(0)
}
