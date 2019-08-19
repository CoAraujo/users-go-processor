package processor

import (
	"encoding/json"
	"github.com/coaraujo/go-processor/domains"
	"github.com/coaraujo/go-processor/infrastructure/config"
	"github.com/coaraujo/go-processor/infrastructure/queue"
	"github.com/coaraujo/go-processor/services/client"
	metauserService "github.com/coaraujo/go-processor/services/metauser"
	"github.com/coaraujo/go-processor/services/oldmetauser"
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
		case msg := <-queue.GetInstance().Notifier(config.UserCreateTopic):
			log.Infof("[Processor Process] Message received. CHANNEL: %s MESSAGE: %s", config.UserCreateTopic, string(msg))
			processUser(msg)

		case msg := <-queue.GetInstance().Notifier(config.UserUpdateTopic):
			log.Infof("[Processor Process] Message received. CHANNEL: %s MESSAGE: %s", config.UserUpdateTopic, string(msg))
			processUser(msg)

		case msg := <-queue.GetInstance().Notifier(config.UserRemovedTopic):
			log.Infof("[Processor Process] Message received. CHANNEL: %s MESSAGE: %s", config.UserRemovedTopic, string(msg))
			processDeletedUser(msg)
		}
	}
}

var processUser = func(msg []byte) {
	//Get message from broker
	var queueResponse domains.QueueResponse
	if err := json.Unmarshal(msg, &queueResponse); err != nil {
		log.Errorf("[Processor processUser] Error to parse. RESPONSE: %s ERROR:", string(msg), err)
		return
	}

	//Get user from client
	log.Infof("[Processor processUser] Processing new MESSAGE: %+v", queueResponse)
	user, err := client.GetInstance().Get(queueResponse.Id)
	if err != nil {
		log.Errorf("[Processor processUser] Error to get user from Client. ERROR: %s", err)
		return
	}

	//Find metauser from mongo
	metauser, err := metauserService.GetInstance().Get(queueResponse.Id)

	//Create new metauser on mongo if it doesnt exist.
	if err == mongo.ErrNoDocuments {
		id, err := metauserService.GetInstance().Insert(user, queueResponse.ClientID)
		if err != nil {
			log.Errorf("[Processor processUser] Error to insert Metauser. ERROR: %s", err)
			return
		}
		log.Infof("[Processor processUser] Message successfully processed. Inserted metauser with ID: %s", id)
		return
	}
	if err != nil {
		log.Errorf("[Processor processUser] Unexpected error to get metauser. ERROR: %s", err)
		return
	}

	//Update metauser on mongo
	if err = metauserService.GetInstance().Update(user, metauser, queueResponse.ClientID); err != nil {
		log.Errorf("[Processor processUser] Error to update metauser on metausers collection. ERROR: %s", err)
		return
	}

	log.Infof("[Processor processUser] Message successfully processed. Updated metauser with ID: %s", metauser.ID)
}

var processDeletedUser = func(msg []byte) {
	//Get message from broker
	var queueResponse domains.QueueResponse
	if err := json.Unmarshal(msg, &queueResponse); err != nil {
		log.Errorf("[Processor processDeletedUser] Error to parse. RESPONSE: %s ERROR: %s", string(msg), err)
		return
	}

	//Find metauser from mongo
	metauser, err := metauserService.GetInstance().Get(queueResponse.Id)
	if err != nil {
		log.Errorf("[Processor processDeletedUser] Unexpected error to get metauser. ERROR: %s", err)
		return
	}

	//Insert metauser on old metausers collection
	_, err = oldmetauser.GetInstance().Insert(metauser)
	if err != nil {
		log.Errorf("[Processor processDeletedUser] Error to move metauser to old metauser collection. ERROR: %s", err)
		return
	}

	if err = metauserService.GetInstance().Delete(metauser.ID); err != nil {
		log.Errorf("[Processor processDeletedUser] Unexpected error to delete metauser. ERROR: %s", err)
		return
	}

	log.Infof("[Processor processDeletedUser] Message successfully processed. Deleted metauser with ID: %s", metauser.ID)
}