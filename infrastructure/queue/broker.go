package queue

import (
	"github.com/coaraujo/go-processor/infrastructure/config"
	"github.com/go-stomp/stomp"
	"github.com/labstack/gommon/log"
	"sync"
	"time"
)

var (
	instance Broker
	once     sync.Once
)

type Broker interface {
	Listen(channel string)
	NewConnection() error
	Disconnect()
	Notifier(channel string) chan []byte
}

type brokerImpl struct {
	conn     *stomp.Conn
	notifier map[string]chan []byte
}

func GetInstance() Broker {
	once.Do(func() {
		instance = &brokerImpl{notifier: make(map[string]chan []byte, 0)}
	})
	return instance
}

func (b *brokerImpl) NewConnection() error {
	conn, err := stomp.Dial(config.ActiveMQProtocol, config.ActiveMQAddress,
		stomp.ConnOpt.Login(config.ActiveMQUser, config.ActiveMQPass),
		stomp.ConnOpt.HeartBeat(0*time.Second, 0*time.Second))
	if err != nil {
		return err
	}

	b.conn = conn
	return nil
}

func (b *brokerImpl) Disconnect() {
	log.Infof("[Broker Disconnect] Disconnecting..")
	err := b.conn.Disconnect()
	if err != nil {
		log.Errorf("[Broker Disconnect] Fail to disconnect. Error: %s ", err)
	}
	log.Infof("[Broker Disconnect] Disconnected")
}

func (b *brokerImpl) Notifier(channel string) chan []byte {
	return b.notifier[channel]
}

func (b *brokerImpl) Listen(channel string) {
	notifier := make(chan []byte)
	b.notifier[channel] = notifier

	log.Infof("[Broker Listen] Subscribing on CHANNEL: %s", channel)
	sub, err := b.conn.Subscribe(channel, stomp.AckAuto)
	if err != nil {
		log.Errorf("[Broker Listen] Fail to subscribe. CHANNEL: %s ERROR: %s", string(channel), err)
	}
	log.Infof("[Broker Listen] Subscribed on CHANNEL: %s", channel)

	for {
		msg := <-sub.C
		log.Infof("[Broker Listen] Received new message. CHANNEL: %s MESSAGE: %s", string(channel), string(msg.Body))
		b.notifier[channel] <- msg.Body
	}
}
