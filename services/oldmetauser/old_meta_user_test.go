package oldmetauser

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

func TestOldMetaUsersImpl_Get_Success(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	oldMetauserMock := &domains.MetaUser{}

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("FindOne", mock.Anything, oldMetausersCollection, mock.Anything, mock.AnythingOfType("*domains.MetaUser")).
		Return(nil).
		Once()

	oldMetauser, err := GetInstance().Get("id")
	assert.Nil(t, err)
	assert.Equal(t, oldMetauser, oldMetauserMock)

	mongoMock.AssertExpectations(t)
}

func TestOldMetaUsersImpl_Get_Error(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	mgoErr := errors.New("error")

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("FindOne", mock.Anything, oldMetausersCollection, mock.Anything, mock.AnythingOfType("*domains.MetaUser")).
		Return(mgoErr).
		Once()

	oldMetauser, err := GetInstance().Get("id")
	assert.NotNil(t, err)
	assert.Equal(t, mgoErr, err)
	assert.Nil(t, oldMetauser)

	mongoMock.AssertExpectations(t)
}

func TestOldMetaUsersImpl_Insert_Success(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	updatedAt := time.Now()
	metauserMock := &domains.MetaUser{}
	mockId := "inserted_id"

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("Insert", mock.Anything, oldMetausersCollection, mock.Anything).
		Return(mockId, nil).
		Once()

	id, err := GetInstance().Insert(metauserMock)
	assert.Nil(t, err)
	assert.Equal(t, id, mockId)
	assert.NotEqual(t, updatedAt, metauserMock.UpdatedAt)

	mongoMock.AssertExpectations(t)
}

func TestOldMetaUsersImpl_Insert_Error(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	updatedAt := time.Now()
	metauserMock := &domains.MetaUser{}
	mockid := "inserted_id"
	mgoErr := errors.New("error")

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("Insert", mock.Anything, oldMetausersCollection, mock.Anything).
		Return("", mgoErr).
		Once()

	id, err := GetInstance().Insert(metauserMock)
	assert.Equal(t, err, mgoErr)
	assert.NotEqual(t, id, mockid)
	assert.NotEqual(t, updatedAt, metauserMock.UpdatedAt)

	mongoMock.AssertExpectations(t)
}
