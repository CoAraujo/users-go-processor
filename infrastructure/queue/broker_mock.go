package queue

import (
	"github.com/go-stomp/stomp"
	"github.com/stretchr/testify/mock"
)

//BrokerMock is a mock for activemq
type BrokerMock struct {
	mock.Mock
}

//Initialize is a mock for db Initialize
func (b *BrokerMock) Initialize() error {
	GetInstance()
	instance = b
	return nil
}

//Listen is a mock for listen a channel
func (b *BrokerMock) Listen(channel string) {}

//NewConnection is a mock for a new connection
func (b *BrokerMock) NewConnection() error {
	args := b.Called()
	return args.Error(0)
}

//Disconnect is a mock for Disconnect
func (b *BrokerMock) Disconnect() {}

//Notifier is a mock for notify a subscriber about channel
func (b *BrokerMock) Notifier(channel string) chan *stomp.Message {
	args := b.Called(channel)
	return args.Get(0).(chan *stomp.Message)
}

//AckMessage is a mock for AckMessage
func (b *BrokerMock) AckMessage(message *stomp.Message) {
	return
}

//RedeliveryMessage is a mock for RedeliveryMessage
func (b *BrokerMock) RedeliveryMessage(message *stomp.Message) {
	return
}