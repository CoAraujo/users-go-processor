package processor

import (
	"github.com/coaraujo/go-processor/domains"
	"github.com/coaraujo/go-processor/services/client"
	"github.com/coaraujo/go-processor/services/metauser"
	"github.com/coaraujo/go-processor/services/oldmetauser"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestProcessUser_UnmarshalError(t *testing.T) {
	clientServiceMock := &client.ClientMock{}
	metauserServiceMock := &metauser.MetaUserMock{}

	msg := []byte("hello world")

	_ = clientServiceMock.Initialize()
	_ = metauserServiceMock.Initialize()

	processUser(msg)

	clientServiceMock.AssertNotCalled(t, "Get", mock.Anything)
	metauserServiceMock.AssertNotCalled(t, "Get", mock.Anything)
	metauserServiceMock.AssertNotCalled(t, "Insert", mock.AnythingOfType("*domains.User"), mock.Anything)
	metauserServiceMock.AssertNotCalled(t, "Update", mock.AnythingOfType("*domains.User"),
		mock.AnythingOfType("*domains.MetaUser"), mock.Anything)

	clientServiceMock.AssertExpectations(t)
	metauserServiceMock.AssertExpectations(t)
}

func TestProcessUser_GetUser_Error(t *testing.T) {
	clientServiceMock := &client.ClientMock{}
	metauserServiceMock := &metauser.MetaUserMock{}

	id := "111111-222-3333-45454545-888990000"
	user := &domains.User{ID: id}
	userError := errors.New("user error")
	msg := []byte("{ \"id\":\"" + id + "\", \"enqueuedAt\": \"2019-08-15T18:15:59-03:00\" }")

	_ = clientServiceMock.Initialize()
	_ = metauserServiceMock.Initialize()

	clientServiceMock.On("Get", id).
		Return(user, userError).
		Once()

	processUser(msg)

	metauserServiceMock.AssertNotCalled(t, "Get", mock.Anything)
	metauserServiceMock.AssertNotCalled(t, "Insert", mock.AnythingOfType("*domains.User"), mock.Anything)
	metauserServiceMock.AssertNotCalled(t, "Update", mock.AnythingOfType("*domains.User"),
		mock.AnythingOfType("*domains.MetaUser"), mock.Anything)

	clientServiceMock.AssertExpectations(t)
	metauserServiceMock.AssertExpectations(t)
}

func TestProcessUser_GetMetaUser_Error(t *testing.T) {
	clientServiceMock := &client.ClientMock{}
	metauserServiceMock := &metauser.MetaUserMock{}

	id := "111111-222-3333-45454545-888990000"
	user := &domains.User{ID: id}
	metauser := &domains.MetaUser{ID: id}
	metaUserError := errors.New("metauser error")
	msg := []byte("{ \"id\":\"" + id + "\", \"enqueuedAt\": \"2019-08-15T18:15:59-03:00\" }")

	_ = clientServiceMock.Initialize()
	_ = metauserServiceMock.Initialize()

	clientServiceMock.On("Get", id).
		Return(user, nil).
		Once()

	metauserServiceMock.On("Get", user.ID).
		Return(metauser, metaUserError).
		Once()

	processUser(msg)

	metauserServiceMock.AssertNotCalled(t, "Insert", mock.AnythingOfType("*domains.User"), mock.Anything)
	metauserServiceMock.AssertNotCalled(t, "Update", mock.AnythingOfType("*domains.User"),
		mock.AnythingOfType("*domains.MetaUser"), mock.Anything)

	clientServiceMock.AssertExpectations(t)
	metauserServiceMock.AssertExpectations(t)
}

func TestProcessUser_GetMetaUser_NotFound(t *testing.T) {
	clientServiceMock := &client.ClientMock{}
	metauserServiceMock := &metauser.MetaUserMock{}

	id := "111111-222-3333-45454545-888990000"
	user := &domains.User{ID: id}
	metauser := &domains.MetaUser{ID: id}
	msg := []byte("{ \"id\":\"" + id + "\", \"enqueuedAt\": \"2019-08-15T18:15:59-03:00\" }")

	_ = clientServiceMock.Initialize()
	_ = metauserServiceMock.Initialize()

	clientServiceMock.On("Get", id).
		Return(user, nil).
		Once()

	metauserServiceMock.On("Get", user.ID).
		Return(metauser, mongo.ErrNoDocuments).
		Once()

	metauserServiceMock.On("Insert", user, mock.Anything).
		Return(id, nil).
		Once()

	processUser(msg)

	metauserServiceMock.AssertNotCalled(t, "Update", mock.AnythingOfType("*domains.User"),
		mock.AnythingOfType("*domains.MetaUser"), mock.Anything)

	clientServiceMock.AssertExpectations(t)
	metauserServiceMock.AssertExpectations(t)
}

func TestProcessUser_InsertMetaUser_Error(t *testing.T) {
	clientServiceMock := &client.ClientMock{}
	metauserServiceMock := &metauser.MetaUserMock{}

	id := "111111-222-3333-45454545-888990000"
	user := &domains.User{ID: id}
	metauser := &domains.MetaUser{ID: id}
	msg := []byte("{ \"id\":\"" + id + "\", \"enqueuedAt\": \"2019-08-15T18:15:59-03:00\" }")
	insertError := errors.New("Insert metauser error.")

	_ = clientServiceMock.Initialize()
	_ = metauserServiceMock.Initialize()

	clientServiceMock.On("Get", id).
		Return(user, nil).
		Once()

	metauserServiceMock.On("Get", user.ID).
		Return(metauser, mongo.ErrNoDocuments).
		Once()

	metauserServiceMock.On("Insert", user, mock.Anything).
		Return("", insertError).
		Once()

	processUser(msg)

	metauserServiceMock.AssertNotCalled(t, "Update", mock.AnythingOfType("*domains.User"),
		mock.AnythingOfType("*domains.MetaUser"), mock.Anything)

	clientServiceMock.AssertExpectations(t)
	metauserServiceMock.AssertExpectations(t)
}

func TestProcessUser_UpdateMetaUser_Error(t *testing.T) {
	clientServiceMock := &client.ClientMock{}
	metauserServiceMock := &metauser.MetaUserMock{}

	id := "111111-222-3333-45454545-888990000"
	user := &domains.User{ID: id}
	metauser := &domains.MetaUser{ID: id}
	msg := []byte("{ \"id\":\"" + id + "\", \"enqueuedAt\": \"2019-08-15T18:15:59-03:00\" }")
	updateError := errors.New("update metauser error.")

	_ = clientServiceMock.Initialize()
	_ = metauserServiceMock.Initialize()

	clientServiceMock.On("Get", id).
		Return(user, nil).
		Once()

	metauserServiceMock.On("Get", user.ID).
		Return(metauser, nil).
		Once()

	metauserServiceMock.On("Update", user, metauser, mock.Anything).
		Return(updateError).
		Once()

	processUser(msg)

	metauserServiceMock.AssertNotCalled(t, "Insert", mock.AnythingOfType("*domains.User"), mock.Anything)
	clientServiceMock.AssertExpectations(t)
	metauserServiceMock.AssertExpectations(t)
}

func TestProcessUser_Success(t *testing.T) {
	clientServiceMock := &client.ClientMock{}
	metauserServiceMock := &metauser.MetaUserMock{}

	id := "111111-222-3333-45454545-888990000"
	user := &domains.User{ID: id}
	metauser := &domains.MetaUser{ID: id}
	msg := []byte("{ \"id\":\"" + id + "\", \"enqueuedAt\": \"2019-08-15T18:15:59-03:00\" }")

	_ = clientServiceMock.Initialize()
	_ = metauserServiceMock.Initialize()

	clientServiceMock.On("Get", id).
		Return(user, nil).
		Once()

	metauserServiceMock.On("Get", user.ID).
		Return(metauser, nil).
		Once()

	metauserServiceMock.On("Update", user, metauser, mock.Anything).
		Return(nil).
		Once()

	processUser(msg)

	metauserServiceMock.AssertNotCalled(t, "Insert", mock.AnythingOfType("*domains.User"), mock.Anything)
	clientServiceMock.AssertExpectations(t)
	metauserServiceMock.AssertExpectations(t)
}

func TestProcessDeletedUser_UnmarshalError(t *testing.T) {
	metauserServiceMock := &metauser.MetaUserMock{}
	oldmetauserServiceMock := &oldmetauser.OldMetaUserMock{}

	msg := []byte("hello world")

	_ = metauserServiceMock.Initialize()
	_ = oldmetauserServiceMock.Initialize()

	processDeletedUser(msg)

	metauserServiceMock.AssertNotCalled(t, "Get", mock.Anything)
	oldmetauserServiceMock.AssertNotCalled(t, "Insert", mock.AnythingOfType("*domains.MetaUser"))
	metauserServiceMock.AssertNotCalled(t, "Delete", mock.Anything)

	oldmetauserServiceMock.AssertExpectations(t)
	metauserServiceMock.AssertExpectations(t)
}

func TestProcessDeletedUser_GetMetaUser_Error(t *testing.T) {
	metauserServiceMock := &metauser.MetaUserMock{}
	oldmetauserServiceMock := &oldmetauser.OldMetaUserMock{}

	id := "111111-222-3333-45454545-888990000"
	metauserMock := &domains.MetaUser{ID: id}
	msg := []byte("{ \"id\":\"" + id + "\", \"enqueuedAt\": \"2019-08-15T18:15:59-03:00\" }")
	getError := errors.New("get metauser error")

	_ = metauserServiceMock.Initialize()
	_ = oldmetauserServiceMock.Initialize()

	metauserServiceMock.On("Get", id).
		Return(metauserMock, getError).
		Once()

	processDeletedUser(msg)

	oldmetauserServiceMock.AssertNotCalled(t, "Insert", mock.AnythingOfType("*domains.MetaUser"))
	metauserServiceMock.AssertNotCalled(t, "Delete", mock.Anything)

	oldmetauserServiceMock.AssertExpectations(t)
	metauserServiceMock.AssertExpectations(t)
}

func TestProcessDeletedUser_InsertOldMetaUser_Error(t *testing.T) {
	metauserServiceMock := &metauser.MetaUserMock{}
	oldmetauserServiceMock := &oldmetauser.OldMetaUserMock{}

	id := "111111-222-3333-45454545-888990000"
	metauserMock := &domains.MetaUser{ID: id}
	msg := []byte("{ \"id\":\"" + id + "\", \"enqueuedAt\": \"2019-08-15T18:15:59-03:00\" }")
	insertError := errors.New("insert metauser error")

	_ = metauserServiceMock.Initialize()
	_ = oldmetauserServiceMock.Initialize()

	metauserServiceMock.On("Get", id).
		Return(metauserMock, nil).
		Once()

	oldmetauserServiceMock.On("Insert", metauserMock).
		Return("id", insertError).
		Once()

	processDeletedUser(msg)

	metauserServiceMock.AssertNotCalled(t, "Delete", mock.Anything)

	oldmetauserServiceMock.AssertExpectations(t)
	metauserServiceMock.AssertExpectations(t)
}

func TestProcessDeletedUser_DeleteMetaUser_Error(t *testing.T) {
	metauserServiceMock := &metauser.MetaUserMock{}
	oldmetauserServiceMock := &oldmetauser.OldMetaUserMock{}

	id := "111111-222-3333-45454545-888990000"
	insertedId := "66666-4444-88888-252525225-6661112222"
	metauserMock := &domains.MetaUser{ID: id}
	msg := []byte("{ \"id\":\"" + id + "\", \"enqueuedAt\": \"2019-08-15T18:15:59-03:00\" }")
	deleteError := errors.New("delete metauser error")

	_ = metauserServiceMock.Initialize()
	_ = oldmetauserServiceMock.Initialize()

	metauserServiceMock.On("Get", id).
		Return(metauserMock, nil).
		Once()

	oldmetauserServiceMock.On("Insert", metauserMock).
		Return(insertedId, nil).
		Once()

	metauserServiceMock.On("Delete", metauserMock.ID).
		Return(deleteError).
		Once()

	processDeletedUser(msg)

	oldmetauserServiceMock.AssertExpectations(t)
	metauserServiceMock.AssertExpectations(t)
}

func TestProcessDeletedUser_Success(t *testing.T) {
	metauserServiceMock := &metauser.MetaUserMock{}
	oldmetauserServiceMock := &oldmetauser.OldMetaUserMock{}

	id := "111111-222-3333-45454545-888990000"
	insertedId := "66666-4444-88888-252525225-6661112222"
	metauserMock := &domains.MetaUser{ID: id}
	msg := []byte("{ \"id\":\"" + id + "\", \"enqueuedAt\": \"2019-08-15T18:15:59-03:00\" }")

	_ = metauserServiceMock.Initialize()
	_ = oldmetauserServiceMock.Initialize()

	metauserServiceMock.On("Get", id).
		Return(metauserMock, nil).
		Once()

	oldmetauserServiceMock.On("Insert", metauserMock).
		Return(insertedId, nil).
		Once()

	metauserServiceMock.On("Delete", metauserMock.ID).
		Return(nil).
		Once()

	processDeletedUser(msg)

	oldmetauserServiceMock.AssertExpectations(t)
	metauserServiceMock.AssertExpectations(t)
}
