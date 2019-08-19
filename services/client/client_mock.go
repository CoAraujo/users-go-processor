package client

import (
	"github.com/coaraujo/go-processor/domains"
	"github.com/stretchr/testify/mock"
)

//ClientMock is a mock for Client
type ClientMock struct {
	mock.Mock
}

//Initialize is a mock for Initialize
func (c *ClientMock) Initialize() error {
	GetInstance()
	instance = c
	return nil
}

//Get is a mock for Get
func (c *ClientMock) Get(id string) (*domains.User, error) {
	args := c.Called(id)
	return args.Get(0).(*domains.User), args.Error(1)
}
