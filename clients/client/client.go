package client

import (
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	once     sync.Once
	instance Client
)

type Client interface {
	Request(method string, url string, body io.Reader) (*http.Response, error)
}

type clientImpl struct{}

func GetInstance() Client {
	once.Do(func() {
		instance = &clientImpl{}
	})
	return instance
}

// Request is a method that make a request to a client
func (c *clientImpl) Request(method string, url string, body io.Reader) (*http.Response, error) {
	request, _ := http.NewRequest(method, url, body)
	request.Header.Add("Origin", "go-processor")

	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1000,
		},
		Timeout: time.Duration(5*time.Second) * time.Second,
	}

	return httpClient.Do(request)
}
