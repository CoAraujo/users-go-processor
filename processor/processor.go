package processor

import (
	"github.com/go-stomp/stomp"
	"time"
	"encoding/json"
	"github.com/coaraujo/users-go-processor/domains"
	"github.com/coaraujo/users-go-processor/infrastructure/config"
	"github.com/coaraujo/users-go-processor/infrastructure/queue"
	"github.com/coaraujo/users-go-processor/services/olduser"
	userService "github.com/coaraujo/users-go-processor/services/user"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

var (
	instance Processor
	once     sync.Once
)

type Processor interface {
	Process()
}

type processorImpl struct {
}

func GetInstance() Processor {
	once.Do(func() {
		instance = &processorImpl{}
	})
	return instance
}

func (p *processorImpl) Process() {
	for {
		select {
		case <-time.After(time.Second * 5):
			continue
			
		case msg := <-queue.GetInstance().Notifier(config.UserCreateTopic):
			log.Infof("[Processor Process] Message received. CHANNEL: %s MESSAGE: %s", config.UserCreateTopic, string(msg.Body))
			processUser(msg)
			continue

		case msg := <-queue.GetInstance().Notifier(config.UserRemovedTopic):
			log.Infof("[Processor Process] Message received. CHANNEL: %s MESSAGE: %s", config.UserRemovedTopic, string(msg.Body))
			processDeletedUser(msg)
			continue
		}
	}
}

var processUser = func(msg *stomp.Message) {
	//Get message from broker
	var user domains.User
	if err := json.Unmarshal(msg.Body, &user); err != nil {
		log.Errorf("[Processor processUser] Error to parse. RESPONSE: %s ERROR:", string(msg.Body), err)
		queue.GetInstance().AckMessage(msg)
		return
	}
	log.Infof("[Processor processUser] Processing new MESSAGE: %+v", user)

	//Find user from mongo
	mongoUser, err := userService.GetInstance().Get(user.ID)

	//Create new user on mongo if it doesnt exist.
	if err == mongo.ErrNoDocuments {
		id, err := userService.GetInstance().Insert(&user)
		if err != nil {
			log.Errorf("[Processor processUser] Error to insert user. ERROR: %s", err)
			queue.GetInstance().RedeliveryMessage(msg)
			return
		}
		queue.GetInstance().AckMessage(msg)
		log.Infof("[Processor processUser] Message successfully processed. Inserted user with ID: %s", id)
		return
	}
	if err != nil {
		log.Errorf("[Processor processUser] Unexpected error to get user. ERROR: %s", err)
		queue.GetInstance().RedeliveryMessage(msg)
		return
	}

	//Update user on mongo
	if err = userService.GetInstance().Update(&user, mongoUser); err != nil {
		log.Errorf("[Processor processUser] Error to update user on users collection. ERROR: %s", err)
		queue.GetInstance().RedeliveryMessage(msg)
		return
	}

	log.Infof("[Processor processUser] Message successfully processed. Updated user with ID: %s", mongoUser.ID)
	queue.GetInstance().AckMessage(msg)
	return
}

var processDeletedUser = func(msg *stomp.Message) {
	//Get message from broker
	var queueResponse domains.User
	if err := json.Unmarshal(msg.Body, &queueResponse); err != nil {
		log.Errorf("[Processor processDeletedUser] Error to parse. RESPONSE: %s ERROR: %s", string(msg.Body), err)
		queue.GetInstance().AckMessage(msg)
		return
	}

	//Find user from mongo
	user, err := userService.GetInstance().Get(queueResponse.ID)
	if err != nil {
		log.Errorf("[Processor processDeletedUser] Unexpected error to get user. ERROR: %s", err)
		queue.GetInstance().RedeliveryMessage(msg)
		return
	}

	//Insert user on old users collection
	_, err = olduser.GetInstance().Insert(user)
	if err != nil {
		log.Errorf("[Processor processDeletedUser] Error to move user to old user collection. ERROR: %s", err)
		queue.GetInstance().RedeliveryMessage(msg)
		return
	}

	if err = userService.GetInstance().Delete(user.ID); err != nil {
		log.Errorf("[Processor processDeletedUser] Unexpected error to delete user. ERROR: %s", err)
		queue.GetInstance().RedeliveryMessage(msg)
		return
	}
	
	log.Infof("[Processor processDeletedUser] Message successfully processed. Deleted user with ID: %s", user.ID)
	queue.GetInstance().AckMessage(msg)
	return
}
