package config

import "os"

const (
	UserUpdateTopic     = "Consumer.Processor.VirtualTopic.user-update"
	UserCreateTopic     = "Consumer.Processor.VirtualTopic.user-create"
	UserRemovedTopic    = "Consumer.Processor.VirtualTopic.user-remove"
)

var (
	ActiveMQAddress  = os.Getenv("ACTIVEMQ_ADDRESS")
	ActiveMQUser     = os.Getenv("ACTIVEMQ_USER")
	ActiveMQPass     = os.Getenv("ACTIVEMQ_PASSWORD")
	ActiveMQProtocol = os.Getenv("ACTIVEMQ_PROTOCOL")

	MongodbAuth     = os.Getenv("MONGODB_AUTH")
	MongodbDatabase = os.Getenv("MONGODB")
	MongodbUser     = os.Getenv("MONGODB_USER")
	MongodbPassword = os.Getenv("MONGODB_PASSWORD")
	MongodbHost     = os.Getenv("MONGODB_HOSTS")
	MongodbPort     = os.Getenv("MONGODB_PORT")

	HttpProtocol = os.Getenv("HTTP_PROTOCOL")
	Host         = os.Getenv("HOST")
)
