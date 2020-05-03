package queue

import (
	"math/rand"
	"github.com/go-stomp/stomp/frame"
	"strconv"
	"github.com/coaraujo/users-go-processor/infrastructure/config"
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
	Notifier(channel string) chan *stomp.Message
	AckMessage(message *stomp.Message)
	RedeliveryMessage(message *stomp.Message)
}

type brokerImpl struct {
	conn     *stomp.Conn
	notifier map[string]chan *stomp.Message
}

func GetInstance() Broker {
	once.Do(func() {
		instance = &brokerImpl{notifier: make(map[string]chan *stomp.Message, 0)}
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

func (b *brokerImpl) AckMessage(message *stomp.Message) {
	if message.ShouldAck() {
		b.conn.Ack(message)
	}
}

func (b *brokerImpl) RedeliveryMessage(message *stomp.Message) {
	log.Infof("[Broker RedeliveryMessage] Redelivering Message: %s", string(message.Body))

	attempt := 1
	attemptsHeader := message.Header.Get("attempts")
	if attemptsHeader != "" {
		attempt, _ = strconv.Atoi(attemptsHeader)
		attempt++
	}

	if attempt > config.MaximumRedeliveries {
		log.Infof("[Broker RedeliveryMessage] Attempt: %d Nack Message: %s", attempt, string(message.Body))
		b.conn.Nack(message)
		return
	}

	log.Infof("[Broker RedeliveryMessage] Resending message. Attempt: %d Message: %s", attempt, string(message.Body))
	b.conn.Ack(message)

	//Redelivery with delay
	go func() {
		time.Sleep(time.Duration(config.RedeliveryDelay) * time.Millisecond)
		b.conn.Send(message.Destination, message.ContentType, message.Body, attemptFunc(attempt))
	}()
}

var attemptFunc = func(attempt int) func(f *frame.Frame) error {
	return func(f *frame.Frame) error {
		f.Header.Add("attempts", strconv.Itoa(attempt))
		return nil
	}
}

func (b *brokerImpl) Disconnect() {
	log.Infof("[Broker Disconnect] Disconnecting..")
	err := b.conn.Disconnect()
	if err != nil {
		log.Errorf("[Broker Disconnect] Fail to disconnect. Error: %s ", err)
	}
	log.Infof("[Broker Disconnect] Disconnected")
}

func (b *brokerImpl) Notifier(channel string) chan *stomp.Message {
	return b.notifier[channel]
}

func (b *brokerImpl) Listen(channel string) {
	notifier := make(chan *stomp.Message)
	b.notifier[channel] = notifier

	log.Infof("[Broker Listen] Subscribing on CHANNEL: %s", channel)
	subID := channel + "-" + strconv.Itoa(rand.Intn(1000))
	sub, err := b.conn.Subscribe(channel, stomp.AckClientIndividual, stomp.SubscribeOpt.Id(subID))
	if err != nil {
		log.Errorf("[Broker Listen] Fail to subscribe. CHANNEL: %s ERROR: %s", string(channel), err)
	}
	log.Infof("[Broker Listen] Subscribed on CHANNEL: %s", channel)

	for {
		msg := <-sub.C
		log.Infof("[Broker Listen] Received new message. CHANNEL: %s MESSAGE: %s", string(channel), string(msg.Body))
		b.notifier[channel] <- msg
	}
}
