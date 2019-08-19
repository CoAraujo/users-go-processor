package client

import (
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
)

//ClientMock is a mock for a client
type ClientMock struct {
	mock.Mock
}

//Initialize is a mock for Initialize
func (c *ClientMock) Initialize() error {
	GetInstance()
	instance = c
	return nil
}

//Request is a mock for Request
func (c *ClientMock) Request(method string, url string, body io.Reader) (*http.Response, error) {
	args := c.Called(method, url, body)
	return args.Get(0).(*http.Response), args.Error(1)
}
