package processor

import (
	"github.com/coaraujo/go-processor/infrastructure/queue"
	"github.com/go-stomp/stomp"
	"github.com/coaraujo/go-processor/domains"
	"github.com/coaraujo/go-processor/services/olduser"
	"github.com/coaraujo/go-processor/services/user"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestProcessUser_UnmarshalError(t *testing.T) {
	userServiceMock := &user.UserMock{}

	msg := &stomp.Message{Body: []byte("hello world")}

	_ = userServiceMock.Initialize()

	processUser(msg)

	userServiceMock.AssertNotCalled(t, "Get", mock.Anything)
	userServiceMock.AssertNotCalled(t, "Insert", mock.AnythingOfType("*domains.User"), mock.Anything)
	userServiceMock.AssertNotCalled(t, "Update", mock.AnythingOfType("*domains.User"),
		mock.AnythingOfType("*domains.User"), mock.Anything)

	userServiceMock.AssertExpectations(t)
}

func TestProcessUser_GetUser_Error(t *testing.T) {
	userServiceMock := &user.UserMock{}
	brokerServiceMock := &queue.BrokerMock{}

	id := "111111-222-3333-45454545-888990000"
	user := &domains.User{ID: id}
	userError := errors.New("user error")
	msg := &stomp.Message{Body: []byte("{ \"_id\":\"" + id + "\", \"enqueuedAt\": \"2019-08-15T18:15:59-03:00\" }")}

	_ = userServiceMock.Initialize()
	_ = brokerServiceMock.Initialize()

	userServiceMock.On("Get", user.ID).
		Return(user, userError).
		Once()

	processUser(msg)

	userServiceMock.AssertNotCalled(t, "Insert", mock.AnythingOfType("*domains.User"), mock.Anything)
	userServiceMock.AssertNotCalled(t, "Update", mock.AnythingOfType("*domains.User"),
		mock.AnythingOfType("*domains.User"), mock.Anything)

	userServiceMock.AssertExpectations(t)
}

func TestProcessUser_GetUser_NotFound(t *testing.T) {
	userServiceMock := &user.UserMock{}
	brokerServiceMock := &queue.BrokerMock{}

	id := "111111-222-3333-45454545-888990000"
	user := &domains.User{ID: id}
	msg := &stomp.Message{Body: []byte("{ \"_id\":\"" + id + "\", \"enqueuedAt\": \"2019-08-15T18:15:59-03:00\" }")}

	_ = userServiceMock.Initialize()
	_ = brokerServiceMock.Initialize()

	userServiceMock.On("Get", user.ID).
		Return(user, mongo.ErrNoDocuments).
		Once()

	userServiceMock.On("Insert", user, mock.Anything).
		Return(id, nil).
		Once()

	processUser(msg)

	userServiceMock.AssertNotCalled(t, "Update", mock.AnythingOfType("*domains.User"),
		mock.AnythingOfType("*domains.User"), mock.Anything)

	userServiceMock.AssertExpectations(t)
}

func TestProcessUser_InsertUser_Error(t *testing.T) {
	userServiceMock := &user.UserMock{}
	brokerServiceMock := &queue.BrokerMock{}

	id := "111111-222-3333-45454545-888990000"
	user := &domains.User{ID: id}
	msg := &stomp.Message{Body: []byte("{ \"_id\":\"" + id + "\", \"enqueuedAt\": \"2019-08-15T18:15:59-03:00\" }")}
	insertError := errors.New("Insert user error")

	_ = userServiceMock.Initialize()
	_ = brokerServiceMock.Initialize()

	userServiceMock.On("Get", user.ID).
		Return(user, mongo.ErrNoDocuments).
		Once()

	userServiceMock.On("Insert", user, mock.Anything).
		Return("", insertError).
		Once()

	processUser(msg)

	userServiceMock.AssertNotCalled(t, "Update", mock.AnythingOfType("*domains.User"),
		mock.AnythingOfType("*domains.User"), mock.Anything)

	userServiceMock.AssertExpectations(t)
}

func TestProcessUser_UpdateUser_Error(t *testing.T) {
	userServiceMock := &user.UserMock{}
	brokerServiceMock := &queue.BrokerMock{}

	id := "111111-222-3333-45454545-888990000"
	newUser := &domains.User{ID: id}
	oldUser := &domains.User{ID: id}
	msg := &stomp.Message{Body: []byte("{ \"_id\":\"" + id + "\", \"enqueuedAt\": \"2019-08-15T18:15:59-03:00\" }")}
	updateError := errors.New("update oldUser error")

	_ = userServiceMock.Initialize()
	_ = brokerServiceMock.Initialize()

	userServiceMock.On("Get", newUser.ID).
		Return(oldUser, nil).
		Once()

	userServiceMock.On("Update", newUser, oldUser, mock.Anything).
		Return(updateError).
		Once()

	processUser(msg)

	userServiceMock.AssertNotCalled(t, "Insert", mock.AnythingOfType("*domains.User"), mock.Anything)
	userServiceMock.AssertExpectations(t)
}

func TestProcessUser_Success(t *testing.T) {
	userServiceMock := &user.UserMock{}
	brokerServiceMock := &queue.BrokerMock{}

	id := "111111-222-3333-45454545-888990000"
	newUser := &domains.User{ID: id}
	oldUser := &domains.User{ID: id}
	msg := &stomp.Message{Body: []byte("{ \"_id\":\"" + id + "\", \"enqueuedAt\": \"2019-08-15T18:15:59-03:00\" }")}

	_ = userServiceMock.Initialize()
	_ = brokerServiceMock.Initialize()

	userServiceMock.On("Get", newUser.ID).
		Return(oldUser, nil).
		Once()

	userServiceMock.On("Update", newUser, oldUser, mock.Anything).
		Return(nil).
		Once()

	processUser(msg)

	userServiceMock.AssertNotCalled(t, "Insert", mock.AnythingOfType("*domains.User"), mock.Anything)
	userServiceMock.AssertExpectations(t)
}

func TestProcessDeletedUser_UnmarshalError(t *testing.T) {
	userServiceMock := &user.UserMock{}
	olduserServiceMock := &olduser.OldUserMock{}
	brokerServiceMock := &queue.BrokerMock{}

	msg := &stomp.Message{Body: []byte("hello world")}

	_ = userServiceMock.Initialize()
	_ = olduserServiceMock.Initialize()
	_ = brokerServiceMock.Initialize()

	processDeletedUser(msg)

	userServiceMock.AssertNotCalled(t, "Get", mock.Anything)
	olduserServiceMock.AssertNotCalled(t, "Insert", mock.AnythingOfType("*domains.User"))
	userServiceMock.AssertNotCalled(t, "Delete", mock.Anything)

	olduserServiceMock.AssertExpectations(t)
	userServiceMock.AssertExpectations(t)
}

func TestProcessDeletedUser_GetUser_Error(t *testing.T) {
	userServiceMock := &user.UserMock{}
	olduserServiceMock := &olduser.OldUserMock{}
	brokerServiceMock := &queue.BrokerMock{}

	id := "111111-222-3333-45454545-888990000"
	userMock := &domains.User{ID: id}
	msg := &stomp.Message{Body: []byte("{ \"_id\":\"" + id + "\", \"enqueuedAt\": \"2019-08-15T18:15:59-03:00\" }")}
	getError := errors.New("get user error")

	_ = userServiceMock.Initialize()
	_ = olduserServiceMock.Initialize()
	_ = brokerServiceMock.Initialize()

	userServiceMock.On("Get", id).
		Return(userMock, getError).
		Once()

	processDeletedUser(msg)

	olduserServiceMock.AssertNotCalled(t, "Insert", mock.AnythingOfType("*domains.User"))
	userServiceMock.AssertNotCalled(t, "Delete", mock.Anything)

	olduserServiceMock.AssertExpectations(t)
	userServiceMock.AssertExpectations(t)
}

func TestProcessDeletedUser_InsertOldUser_Error(t *testing.T) {
	userServiceMock := &user.UserMock{}
	olduserServiceMock := &olduser.OldUserMock{}
	brokerServiceMock := &queue.BrokerMock{}

	id := "111111-222-3333-45454545-888990000"
	userMock := &domains.User{ID: id}
	msg := &stomp.Message{Body: []byte("{ \"_id\":\"" + id + "\", \"enqueuedAt\": \"2019-08-15T18:15:59-03:00\" }")}
	insertError := errors.New("insert user error")

	_ = userServiceMock.Initialize()
	_ = olduserServiceMock.Initialize()
	_ = brokerServiceMock.Initialize()

	userServiceMock.On("Get", id).
		Return(userMock, nil).
		Once()

	olduserServiceMock.On("Insert", userMock).
		Return("id", insertError).
		Once()

	processDeletedUser(msg)

	userServiceMock.AssertNotCalled(t, "Delete", mock.Anything)

	olduserServiceMock.AssertExpectations(t)
	userServiceMock.AssertExpectations(t)
}

func TestProcessDeletedUser_DeleteUser_Error(t *testing.T) {
	userServiceMock := &user.UserMock{}
	olduserServiceMock := &olduser.OldUserMock{}
	brokerServiceMock := &queue.BrokerMock{}

	id := "111111-222-3333-45454545-888990000"
	insertedID := "66666-4444-88888-252525225-6661112222"
	userMock := &domains.User{ID: id}
	msg := &stomp.Message{Body: []byte("{ \"_id\":\"" + id + "\", \"enqueuedAt\": \"2019-08-15T18:15:59-03:00\" }")}
	deleteError := errors.New("delete user error")

	_ = userServiceMock.Initialize()
	_ = olduserServiceMock.Initialize()
	_ = brokerServiceMock.Initialize()

	userServiceMock.On("Get", id).
		Return(userMock, nil).
		Once()

	olduserServiceMock.On("Insert", userMock).
		Return(insertedID, nil).
		Once()

	userServiceMock.On("Delete", userMock.ID).
		Return(deleteError).
		Once()

	processDeletedUser(msg)

	olduserServiceMock.AssertExpectations(t)
	userServiceMock.AssertExpectations(t)
}

func TestProcessDeletedUser_Success(t *testing.T) {
	userServiceMock := &user.UserMock{}
	olduserServiceMock := &olduser.OldUserMock{}
	brokerServiceMock := &queue.BrokerMock{}

	id := "111111-222-3333-45454545-888990000"
	insertedID := "66666-4444-88888-252525225-6661112222"
	userMock := &domains.User{ID: id}
	msg := &stomp.Message{Body: []byte("{ \"_id\":\"" + id + "\", \"enqueuedAt\": \"2019-08-15T18:15:59-03:00\" }")}

	_ = userServiceMock.Initialize()
	_ = olduserServiceMock.Initialize()
	_ = brokerServiceMock.Initialize()

	userServiceMock.On("Get", id).
		Return(userMock, nil).
		Once()

	olduserServiceMock.On("Insert", userMock).
		Return(insertedID, nil).
		Once()

	userServiceMock.On("Delete", userMock.ID).
		Return(nil).
		Once()

	processDeletedUser(msg)

	olduserServiceMock.AssertExpectations(t)
	userServiceMock.AssertExpectations(t)
}
