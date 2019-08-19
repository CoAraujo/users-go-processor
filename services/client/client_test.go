package client

import (
	"github.com/coaraujo/go-processor/clients/client"
	"github.com/coaraujo/go-processor/domains"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestClientImpl_Get_Error(t *testing.T) {
	clientMock := &client.ClientMock{}
	httpResponse := &http.Response{}
	gErr := errors.New("error")

	_ = clientMock.Initialize()
	clientMock.On("Request", mock.Anything, mock.Anything, mock.Anything).
		Return(httpResponse, gErr).
		Once()

	user, err := GetInstance().Get("id")
	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.Equal(t, gErr, err)

	clientMock.AssertExpectations(t)
}

func TestClientImpl_Get_Error_StatusNotOK(t *testing.T) {
	clientMock := &client.ClientMock{}
	body := ioutil.NopCloser(strings.NewReader("body"))
	httpResponse := &http.Response{StatusCode: http.StatusBadGateway, Body: body}

	_ = clientMock.Initialize()
	clientMock.On("Request", http.MethodGet, mock.Anything, mock.Anything).
		Return(httpResponse, nil).
		Once()

	user, err := GetInstance().Get("id")
	assert.Nil(t, user)
	assert.NotNil(t, err)

	clientMock.AssertExpectations(t)
}

func TestClientImpl_Get_Error_CannotDecode(t *testing.T) {
	clientMock := &client.ClientMock{}
	body := ioutil.NopCloser(strings.NewReader("body"))
	httpResponse := &http.Response{StatusCode: http.StatusOK, Body: body}

	_ = clientMock.Initialize()
	clientMock.On("Request", http.MethodGet, mock.Anything, mock.Anything).
		Return(httpResponse, nil).
		Once()

	user, err := GetInstance().Get("id")
	assert.Nil(t, user)
	assert.NotNil(t, err)

	clientMock.AssertExpectations(t)
}

func TestClientImpl_Get_Success(t *testing.T) {
	clientMock := &client.ClientMock{}
	body := ioutil.NopCloser(strings.NewReader("{\"id\":\"123\"}"))
	httpResponse := &http.Response{StatusCode: http.StatusOK, Body: body}

	expectedUser := domains.User{ID: "123"}

	_ = clientMock.Initialize()
	clientMock.On("Request", http.MethodGet, mock.Anything, mock.Anything).
		Return(httpResponse, nil).
		Once()

	user, err := GetInstance().Get("id")
	assert.NotNil(t, user)
	assert.Equal(t, user.ID, expectedUser.ID)
	assert.Nil(t, err)

	clientMock.AssertExpectations(t)
}
