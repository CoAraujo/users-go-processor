package olduser

import (
	"context"
	"github.com/coaraujo/go-processor/domains"
	"github.com/coaraujo/go-processor/infrastructure/storage"
	"sync"
	"time"
)

const (
	oldUsersCollection = "old_users"
)

var (
	instance OldUsers
	once     sync.Once
)

type OldUsers interface {
	Get(id string) (*domains.User, error)
	Insert(user *domains.User) (string, error)
}

type oldUsersImpl struct{}

func GetInstance() OldUsers {
	once.Do(func() {
		instance = &oldUsersImpl{}
	})
	return instance
}

func (o *oldUsersImpl) Get(id string) (*domains.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var user domains.User
	if mgoErr := storage.GetInstance().FindOne(ctx, oldUsersCollection, map[string]interface{}{"_id": id}, &user); mgoErr != nil {
		return nil, mgoErr
	}

	return &user, nil
}

func (o *oldUsersImpl) Insert(user *domains.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	user.UpdatedAt = time.Now()

	id, mgoErr := storage.GetInstance().Insert(ctx, oldUsersCollection, user)
	if mgoErr != nil {
		return "", mgoErr
	}

	return id.(string), nil
}
