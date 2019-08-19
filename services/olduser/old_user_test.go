package olduser

import (
	"context"
	"github.com/coaraujo/go-processor/domains"
	"github.com/coaraujo/go-processor/infrastructure/storage"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestOldUsersImpl_Get_Success(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	oldUserMock := &domains.User{}

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("FindOne", mock.Anything, oldUsersCollection, mock.Anything, mock.AnythingOfType("*domains.User")).
		Return(nil).
		Once()

	oldUser, err := GetInstance().Get("id")
	assert.Nil(t, err)
	assert.Equal(t, oldUser, oldUserMock)

	mongoMock.AssertExpectations(t)
}

func TestOldUsersImpl_Get_Error(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	mgoErr := errors.New("error")

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("FindOne", mock.Anything, oldUsersCollection, mock.Anything, mock.AnythingOfType("*domains.User")).
		Return(mgoErr).
		Once()

	oldUser, err := GetInstance().Get("id")
	assert.NotNil(t, err)
	assert.Equal(t, mgoErr, err)
	assert.Nil(t, oldUser)

	mongoMock.AssertExpectations(t)
}

func TestOldUsersImpl_Insert_Success(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	updatedAt := time.Now()
	oldUserMock := &domains.User{}
	mockId := "inserted_id"

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("Insert", mock.Anything, oldUsersCollection, mock.Anything).
		Return(mockId, nil).
		Once()

	id, err := GetInstance().Insert(oldUserMock)
	assert.Nil(t, err)
	assert.Equal(t, id, mockId)
	assert.NotEqual(t, updatedAt, oldUserMock.UpdatedAt)

	mongoMock.AssertExpectations(t)
}

func TestOldUsersImpl_Insert_Error(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	updatedAt := time.Now()
	userMock := &domains.User{}
	mockid := "inserted_id"
	mgoErr := errors.New("error")

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("Insert", mock.Anything, oldUsersCollection, mock.Anything).
		Return("", mgoErr).
		Once()

	id, err := GetInstance().Insert(userMock)
	assert.Equal(t, err, mgoErr)
	assert.NotEqual(t, id, mockid)
	assert.NotEqual(t, updatedAt, userMock.UpdatedAt)

	mongoMock.AssertExpectations(t)
}
