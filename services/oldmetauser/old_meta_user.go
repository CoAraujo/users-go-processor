package oldmetauser

import (
	"context"
	"github.com/coaraujo/go-processor/domains"
	"github.com/coaraujo/go-processor/infrastructure/storage"
	"sync"
	"time"
)

const (
	oldMetausersCollection = "old_metausers"
)

var (
	instance OldMetaUsers
	once     sync.Once
)

type OldMetaUsers interface {
	Get(id string) (*domains.MetaUser, error)
	Insert(metauser *domains.MetaUser) (string, error)
}

type oldMetaUsersImpl struct{}

func GetInstance() OldMetaUsers {
	once.Do(func() {
		instance = &oldMetaUsersImpl{}
	})
	return instance
}

func (m *oldMetaUsersImpl) Get(id string) (*domains.MetaUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var metauser domains.MetaUser
	if mgoErr := storage.GetInstance().FindOne(ctx, oldMetausersCollection, map[string]interface{}{"_id": id}, &metauser); mgoErr != nil {
		return nil, mgoErr
	}

	return &metauser, nil
}

func (m *oldMetaUsersImpl) Insert(metauser *domains.MetaUser) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	metauser.UpdatedAt = time.Now()

	id, mgoErr := storage.GetInstance().Insert(ctx, oldMetausersCollection, metauser)
	if mgoErr != nil {
		return "", mgoErr
	}

	return id.(string), nil
}
