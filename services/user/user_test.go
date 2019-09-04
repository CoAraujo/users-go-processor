package user

import (
	"go.mongodb.org/mongo-driver/mongo"
	"context"
	"errors"
	"github.com/coaraujo/go-processor/domains"
	"github.com/coaraujo/go-processor/infrastructure/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestUsersImpl_Get_Success(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	userMock := &domains.User{}

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("FindOne", mock.Anything, usersCollection, mock.Anything, mock.AnythingOfType("*domains.User")).
		Return(nil).
		Once()

	user, err := GetInstance().Get("id")
	assert.Nil(t, err)
	assert.Equal(t, user, userMock)

	mongoMock.AssertExpectations(t)
}

func TestUsersImpl_Get_Error(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	mgoErr := errors.New("error")

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("FindOne", mock.Anything, usersCollection, mock.Anything, mock.AnythingOfType("*domains.User")).
		Return(mgoErr).
		Once()

	user, err := GetInstance().Get("id")
	assert.NotNil(t, err)
	assert.Equal(t, mgoErr, err)
	assert.Nil(t, user)

	mongoMock.AssertExpectations(t)
}

func TestUsersImpl_Insert_Success(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	updatedAt := time.Now()
	user := &domains.User{UpdatedAt: updatedAt}

	mockid := "inserted_id"

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("Insert", mock.Anything, usersCollection, mock.Anything).
		Return(mockid, nil).
		Once()

	id, err := GetInstance().Insert(user)
	assert.Nil(t, err)
	assert.Equal(t, id, mockid)

	mongoMock.AssertExpectations(t)
}

func TestUsersImpl_Insert_GenericError(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	updatedAt := time.Now()
	user := &domains.User{UpdatedAt: updatedAt}
	mgoErr := errors.New("error")
	mockid := "inserted_id"

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("Insert", mock.Anything, usersCollection, mock.Anything).
		Return("", mgoErr).
		Once()

	id, err := GetInstance().Insert(user)
	assert.Equal(t, err, mgoErr)
	assert.NotEqual(t, id, mockid)

	mongoMock.AssertExpectations(t)
}

func TestUsersImpl_Insert_UpdatedAtNull_Error(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	user := &domains.User{}
	mockid := "inserted_id"

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("Insert", mock.Anything, usersCollection, mock.Anything).
		Return(mockid, nil).
		Once()

	id, err := GetInstance().Insert(user)
	assert.Nil(t, err)
	assert.Equal(t, id, mockid)
	assert.NotNil(t, user.UpdatedAt)

	mongoMock.AssertExpectations(t)
}

func TestUsersImpl_Delete_Success(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}

	mockId := "Id"

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("Remove", mock.Anything, usersCollection, mock.Anything).
		Return(nil).
		Once()

	err := GetInstance().Delete(mockId)
	assert.Nil(t, err)

	mongoMock.AssertExpectations(t)
}

func TestUsersImpl_Delete_Error(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}

	mockId := "Id"
	mgoErr := errors.New("error")

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("Remove", mock.Anything, usersCollection, mock.Anything).
		Return(mgoErr).
		Once()

	err := GetInstance().Delete(mockId)
	assert.NotNil(t, err)
	assert.Equal(t, err, mgoErr)

	mongoMock.AssertExpectations(t)
}

func TestUsersImpl_Update_Success(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}

	mockClientId := "clientId"
	username := "newUsername"
	email := "oldEmail"
	status := "unchangedStatus"
	phone1 := &domains.Phone{Phone: "phone1", UpdatedAt: time.Now()}
	phone2 := &domains.Phone{Phone: "phone2", UpdatedAt: time.Now()}

	oldUser := domains.User{
		Email:     email,
		Username:  "username",
		Name:      "name1",
		Gender:    "gender1",
		Status:    status,
		BirthDate: "birthdate1",
		Phones:    phone1,
		ClientID:  mockClientId,
	}

	newUser := domains.User{
		Username:  username,
		Name:      "name2",
		Gender:    "gender2",
		Status:    status,
		BirthDate: "birthdate2",
		Phones:    phone2,
		ClientID:  mockClientId,
	}

	expectedUser := domains.User{
		Email:     email,
		Username:  username,
		Name:      "name2",
		Gender:    "gender2",
		Status:    status,
		BirthDate: "birthdate2",
		Phones:    phone2,
		ClientID:  mockClientId,
	}

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("UpdateOne", mock.Anything, usersCollection, mock.Anything, mock.Anything).
		Return(&mongo.UpdateResult{}, nil).
		Once()

	err := GetInstance().Update(&newUser, &oldUser)
	assert.Nil(t, err)
	assert.Equal(t, oldUser, expectedUser)

	mongoMock.AssertExpectations(t)
}

func TestUsersImpl_Update_Error(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}

	mgoErr := errors.New("error")

	user := &domains.User{}
	userMock := &domains.User{}

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("UpdateOne", mock.Anything, usersCollection, mock.Anything, mock.Anything).
		Return(&mongo.UpdateResult{}, mgoErr).
		Once()

	err := GetInstance().Update(user, userMock)
	assert.NotNil(t, err)
	assert.Equal(t, err, mgoErr)

	mongoMock.AssertExpectations(t)
}

func TestUsersImpl_updateNewUserValues_ChangeAllValues(t *testing.T) {
	phone1 := &domains.Phone{Phone: "phone1", UpdatedAt: time.Now()}
	phone2 := &domains.Phone{Phone: "phone2", UpdatedAt: time.Now()}

	oldUser := domains.User{
		Email:     "email1",
		Username:  "username1",
		Name:      "name1",
		Gender:    "gender1",
		Status:    "status1",
		BirthDate: "birthdate1",
		Phones:    phone1,
		ClientID:  "client1",
	}

	newUser := domains.User{
		Email:     "email2",
		Username:  "username2",
		Name:      "name2",
		Gender:    "gender2",
		Status:    "status2",
		BirthDate: "birthdate2",
		Phones:    phone2,
		ClientID:  "client2",
	}

	expectedUser := domains.User{
		Email:     "email2",
		Username:  "username2",
		Name:      "name2",
		Gender:    "gender2",
		Status:    "status2",
		BirthDate: "birthdate2",
		Phones:    phone2,
		ClientID:  "client2",
	}

	updateNewUserValues(&oldUser, &newUser)

	assert.Equal(t, newUser, expectedUser)
	assert.Equal(t, oldUser, expectedUser)
}

func TestUsersImpl_updateNewUserValues_ChangeSomeValues(t *testing.T) {
	phone1 := &domains.Phone{Phone: "phone1", UpdatedAt: time.Now()}
	phone2 := &domains.Phone{Phone: "phone2", UpdatedAt: time.Now()}

	oldUser := domains.User{
		Email:    "email1",
		Phones:   phone1,
		Username: "username1",
	}

	newUser := domains.User{
		Email:  "email1",
		Phones: phone2,
		Status: "status2",
	}

	expectedUser := domains.User{
		Email:    "email1",
		Phones:   phone2,
		Status:   "status2",
		Username: "username1",
	}

	updateNewUserValues(&oldUser, &newUser)

	assert.Equal(t, oldUser, expectedUser)
	assert.NotEqual(t, newUser, expectedUser)
}

func TestUsersImpl_isEqual_True(t *testing.T) {
	phone1 := &domains.Phone{Phone: "phone1", UpdatedAt: time.Now()}
	phone2 := &domains.Phone{Phone: "phone1", UpdatedAt: time.Now()}

	user := domains.User{
		Email:    "email1",
		Phones:   phone1,
		Status:   "status1",
		Username: "username1",
	}

	anotherUser := domains.User{
		Email:    "email2",
		Phones:   phone2,
		Status:   "status2",
		Username: "username2",
	}

	resp := isEqual(*user.Phones, *anotherUser.Phones)
	assert.True(t, resp)
}

func TestUsersImpl_isEqual_False(t *testing.T) {
	phone1 := &domains.Phone{Phone: "phone1", UpdatedAt: time.Now()}
	phone2 := &domains.Phone{Phone: "phone2", UpdatedAt: time.Now()}

	user := domains.User{
		Email:    "email1",
		Phones:   phone1,
		Status:   "status1",
		Username: "username1",
	}

	anotherUser := domains.User{
		Email:    "email2",
		Phones:   phone2,
		Status:   "status2",
		Username: "username2",
	}

	resp := isEqual(*user.Phones, *anotherUser.Phones)
	assert.False(t, resp)
}

func TestUsersImpl_validateUpdatedAt_WhenItExist(t *testing.T) {
	updatedAt := time.Now()
	user := &domains.User{UpdatedAt: updatedAt}

	validateUpdatedAt(user)

	assert.Equal(t, user.UpdatedAt, updatedAt)
}

func TestUsersImpl_validateUpdatedAt_WhenItDoesntExist(t *testing.T) {
	updatedAt := time.Now()
	user := &domains.User{}

	validateUpdatedAt(user)

	assert.NotNil(t, user.UpdatedAt)
	assert.NotEqual(t, user.UpdatedAt, &updatedAt)
}
