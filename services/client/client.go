package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coaraujo/go-processor/clients/client"
	"github.com/coaraujo/go-processor/domains"
	"github.com/coaraujo/go-processor/infrastructure/config"
	"github.com/labstack/gommon/log"
	"net/http"
	"sync"
)

var (
	instance Client
	once     sync.Once
)

type Client interface {
	Get(id string) (*domains.User, error)
}

type clientImpl struct{}

func GetInstance() Client {
	once.Do(func() {
		instance = &clientImpl{}
	})
	return instance
}

func (g *clientImpl) Get(id string) (*domains.User, error) {
	url := fmt.Sprintf("%s://%s/v2/users/%s", config.HttpProtocol, config.Host, id)
	resp, err := client.GetInstance().Request(http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("[Client GetUser] Error to call client to get user. Id: %s Error:", id, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Error when calling client to get user. ID: %s Status: %d", id, resp.StatusCode)
		return nil, errors.New(msg)
	}

	var user domains.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		log.Errorf("[Client GetUser] Error to parse response body. Error: ", err)
		return nil, err
	}

	return &user, nil
}
