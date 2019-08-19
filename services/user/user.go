package user

import (
	"context"
	"github.com/coaraujo/go-processor/domains"
	"github.com/coaraujo/go-processor/infrastructure/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/labstack/gommon/log"
	"sync"
	"time"
)

const (
	usersCollection = "users"
)

var (
	instance Users
	once     sync.Once
)

type Users interface {
	Get(id string) (*domains.User, error)
	Insert(user *domains.User) (string, error)
	Update(newUser *domains.User, oldUser *domains.User) error
	Delete(id string) error
}

type usersImpl struct{}

func GetInstance() Users {
	once.Do(func() {
		instance = &usersImpl{}
	})
	return instance
}

func (u *usersImpl) Get(id string) (*domains.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var user domains.User
	if mgoErr := storage.GetInstance().FindOne(ctx, usersCollection, map[string]interface{}{"_id": id}, &user); mgoErr != nil {
		return nil, mgoErr
	}

	return &user, nil
}

func (u *usersImpl) Insert(user *domains.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	validateUpdatedAt(user)

	id, mgoErr := storage.GetInstance().Insert(ctx, usersCollection, user)
	if mgoErr != nil {
		return id.(string), mgoErr
	}

	return id.(string), nil
}

func (u *usersImpl) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if mgoErr := storage.GetInstance().Remove(ctx, usersCollection, map[string]interface{}{"_id": id}); mgoErr != nil {
		return mgoErr
	}

	return nil
}

func (u *usersImpl) Update(newUser *domains.User, oldUser *domains.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	validateUpdatedAt(newUser)
	updateNewUserValues(oldUser, newUser)

	if mgoErr := storage.GetInstance().Update(ctx, usersCollection, map[string]interface{}{"_id": oldUser.ID},
		map[string]interface{}{"$set": &oldUser}); mgoErr != nil {
		return mgoErr
	}

	return nil
}

var updateNewUserValues = func(oldUser *domains.User, newUser *domains.User) {
	if newUser.Phones != nil && !isEqual(*oldUser.Phones, *newUser.Phones) {
		oldUser.Phones = newUser.Phones
	}
	if newUser.Status != "" {
		oldUser.Status = newUser.Status
	}
	if newUser.Username != "" {
		oldUser.Username = newUser.Username
	}
	if newUser.Email != "" {
		oldUser.Email = newUser.Email
	}
	if newUser.Name != "" {
		oldUser.Name = newUser.Name
	}
	if newUser.ClientID != "" {
		oldUser.ClientID = newUser.ClientID
	}
	if newUser.Gender != "" {
		oldUser.Gender = newUser.Gender
	}
	if newUser.BirthDate != "" {
		oldUser.BirthDate = newUser.BirthDate
	}
	oldUser.UpdatedAt = newUser.UpdatedAt
}

var isEqual = func(interface1, interface2 interface{}) bool {
	return cmp.Equal(interface1, interface2, cmpopts.IgnoreFields(interface1, "UpdatedAt"))
}

var validateUpdatedAt = func(user *domains.User) {
	if &user.UpdatedAt == nil {
		log.Errorf("[User validateUpdatedAt] user without updatedAt. Setting updatedAt with time.Now(). Id: %s ", user.ID)
		updatedAt := time.Now()
		user.UpdatedAt = updatedAt
	}
}
