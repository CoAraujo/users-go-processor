package metauser

import (
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

func TestMetaUsersImpl_Get_Success(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	metauserMock := &domains.MetaUser{}

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("FindOne", mock.Anything, metausersCollection, mock.Anything, mock.AnythingOfType("*domains.MetaUser")).
		Return(nil).
		Once()

	metauser, err := GetInstance().Get("id")
	assert.Nil(t, err)
	assert.Equal(t, metauser, metauserMock)

	mongoMock.AssertExpectations(t)
}

func TestMetaUsersImpl_Get_Error(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	mgoErr := errors.New("error")

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("FindOne", mock.Anything, metausersCollection, mock.Anything, mock.AnythingOfType("*domains.MetaUser")).
		Return(mgoErr).
		Once()

	metauser, err := GetInstance().Get("id")
	assert.NotNil(t, err)
	assert.Equal(t, mgoErr, err)
	assert.Nil(t, metauser)

	mongoMock.AssertExpectations(t)
}

func TestMetaUsersImpl_Insert_Success(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	updatedAt := time.Now()
	mockClientId := "clientId"
	user := &domains.User{UpdatedAt: &updatedAt}

	mockid := "inserted_id"

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("Insert", mock.Anything, metausersCollection, mock.Anything).
		Return(mockid, nil).
		Once()

	id, err := GetInstance().Insert(user, mockClientId)
	assert.Nil(t, err)
	assert.Equal(t, id, mockid)

	mongoMock.AssertExpectations(t)
}

func TestMetaUsersImpl_Insert_GenericError(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	updatedAt := time.Now()
	mockClientId := "clientId"
	user := &domains.User{UpdatedAt: &updatedAt}
	mgoErr := errors.New("error")
	mockid := "inserted_id"

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("Insert", mock.Anything, metausersCollection, mock.Anything).
		Return("", mgoErr).
		Once()

	id, err := GetInstance().Insert(user, mockClientId)
	assert.Equal(t, err, mgoErr)
	assert.NotEqual(t, id, mockid)

	mongoMock.AssertExpectations(t)
}

func TestMetaUsersImpl_Insert_UpdatedAtNull_Error(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}
	mockClientId := "clientId"
	user := &domains.User{}
	mockid := "inserted_id"

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("Insert", mock.Anything, metausersCollection, mock.Anything).
		Return(mockid, nil).
		Once()

	id, err := GetInstance().Insert(user, mockClientId)
	assert.Nil(t, err)
	assert.Equal(t, id, mockid)
	assert.NotNil(t, *user.UpdatedAt)

	mongoMock.AssertExpectations(t)
}

func TestMetaUsersImpl_Delete_Success(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}

	mockId := "Id"

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("Remove", mock.Anything, metausersCollection, mock.Anything).
		Return(nil).
		Once()

	err := GetInstance().Delete(mockId)
	assert.Nil(t, err)

	mongoMock.AssertExpectations(t)
}

func TestMetaUsersImpl_Delete_Error(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}

	mockId := "Id"
	mgoErr := errors.New("error")

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("Remove", mock.Anything, metausersCollection, mock.Anything).
		Return(mgoErr).
		Once()

	err := GetInstance().Delete(mockId)
	assert.NotNil(t, err)
	assert.Equal(t, err, mgoErr)

	mongoMock.AssertExpectations(t)
}

func TestMetaUsersImpl_Update_Success(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}

	mockClientId := "clientId"
	username := "newUsername"
	email := "oldEmail"
	status := "unchangedStatus"

	user := &domains.User{
		Username: username,
		Status:   status,
	}

	metauserMock := &domains.MetaUser{
		Email: &domains.EmailMetadata{
			NewValue: &domains.EmailValueMetadata{Email: email},
		},
		Status: &domains.StatusMetadata{
			NewValue: &domains.StatusValueMetadata{Status: status},
		},
	}

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("Update", mock.Anything, metausersCollection, mock.Anything, mock.Anything).
		Return(nil).
		Once()

	err := GetInstance().Update(user, metauserMock, mockClientId)
	assert.Nil(t, err)
	assert.Equal(t, metauserMock.Username.NewValue.Username, username)
	assert.Equal(t, metauserMock.Username.NewValue.ClientID, mockClientId)
	assert.Equal(t, metauserMock.Email.LastValue.Email, email)
	assert.Equal(t, metauserMock.Email.NewValue.Email, user.Email)
	assert.Equal(t, metauserMock.Email.NewValue.Email, user.Email)
	assert.Equal(t, metauserMock.Status.NewValue.Status, status)
	assert.Nil(t, metauserMock.Status.LastValue)

	mongoMock.AssertExpectations(t)
}

func TestMetaUsersImpl_Update_Error(t *testing.T) {
	mongoMock := &storage.DataAccessLayerMock{}

	mgoErr := errors.New("error")
	mockClientId := "clientId"

	user := &domains.User{}
	metauserMock := &domains.MetaUser{}

	_ = mongoMock.Initialize(context.Background(), options.Credential{}, mock.Anything, mock.Anything)
	mongoMock.On("Update", mock.Anything, metausersCollection, mock.Anything, mock.Anything).
		Return(mgoErr).
		Once()

	err := GetInstance().Update(user, metauserMock, mockClientId)
	assert.NotNil(t, err)
	assert.Equal(t, err, mgoErr)

	mongoMock.AssertExpectations(t)
}

func TestMetaUsersImpl_updateNewMetauserValues_ChangeAllValues(t *testing.T) {
	address1 := &domains.AddressValueMetadata{City: &domains.CityMetadata{ID: 1, Name: "1", IbgeID: 1}, Neighborhood: "neighborhood1"}
	address2 := &domains.AddressValueMetadata{City: &domains.CityMetadata{ID: 2, Name: "2", IbgeID: 2}, Neighborhood: "neighborhood2"}
	phone1 := &domains.PhoneValueMetadata{Phone: "phone1"}
	phone2 := &domains.PhoneValueMetadata{Phone: "phone2"}
	email1 := &domains.EmailValueMetadata{Email: "email1"}
	email2 := &domains.EmailValueMetadata{Email: "email2"}
	status1 := &domains.StatusValueMetadata{Status: "status1"}
	status2 := &domains.StatusValueMetadata{Status: "status2"}
	username1 := &domains.UsernameValueMetadata{Username: "username1"}
	username2 := &domains.UsernameValueMetadata{Username: "username2"}
	password1 := &domains.PasswordValueMetadata{ClientID: "clientid1"}
	password2 := &domains.PasswordValueMetadata{ClientID: "clientid2"}

	metauser := &domains.MetaUser{
		Address:  &domains.AddressMetadata{NewValue: address1},
		Phones:   &domains.PhonesMetadata{NewValue: phone1},
		Email:    &domains.EmailMetadata{NewValue: email1},
		Status:   &domains.StatusMetadata{NewValue: status1},
		Username: &domains.UsernameMetadata{NewValue: username1},
		Password: &domains.PasswordMetadata{NewValue: password1},
	}

	tempMetauser := domains.MetaUser{
		Address:  &domains.AddressMetadata{NewValue: address2},
		Phones:   &domains.PhonesMetadata{NewValue: phone2},
		Email:    &domains.EmailMetadata{NewValue: email2},
		Status:   &domains.StatusMetadata{NewValue: status2},
		Username: &domains.UsernameMetadata{NewValue: username2},
		Password: &domains.PasswordMetadata{NewValue: password2},
	}

	expectedMetauser := &domains.MetaUser{
		Address:  &domains.AddressMetadata{NewValue: address2, LastValue: address1},
		Phones:   &domains.PhonesMetadata{NewValue: phone2, LastValue: phone1},
		Email:    &domains.EmailMetadata{NewValue: email2, LastValue: email1},
		Status:   &domains.StatusMetadata{NewValue: status2, LastValue: status1},
		Username: &domains.UsernameMetadata{NewValue: username2, LastValue: username1},
		Password: &domains.PasswordMetadata{NewValue: password1},
	}

	updateNewMetauserValues(metauser, tempMetauser)

	assert.Equal(t, metauser, expectedMetauser)
}

func TestMetaUsersImpl_updateNewMetauserValues_ChangeSomeValues(t *testing.T) {
	address1 := &domains.AddressValueMetadata{City: &domains.CityMetadata{ID: 1, Name: "1", IbgeID: 1}, Neighborhood: "neighborhood1"}
	address2 := &domains.AddressValueMetadata{City: &domains.CityMetadata{ID: 2, Name: "2", IbgeID: 2}, Neighborhood: "neighborhood2"}
	phone2 := &domains.PhoneValueMetadata{Phone: "phone2"}
	email1 := &domains.EmailValueMetadata{Email: "email1"}
	status1 := &domains.StatusValueMetadata{Status: "status1"}
	password2 := &domains.PasswordValueMetadata{ClientID: "clientid2"}

	metauser := &domains.MetaUser{
		Address: &domains.AddressMetadata{NewValue: address1},
		Email:   &domains.EmailMetadata{NewValue: email1},
		Status:  &domains.StatusMetadata{NewValue: status1},
	}

	tempMetauser := domains.MetaUser{
		Address:  &domains.AddressMetadata{NewValue: address2},
		Phones:   &domains.PhonesMetadata{NewValue: phone2},
		Password: &domains.PasswordMetadata{NewValue: password2},
	}

	expectedMetauser := &domains.MetaUser{
		Address: &domains.AddressMetadata{NewValue: address2, LastValue: address1},
		Phones:  &domains.PhonesMetadata{NewValue: phone2},
		Email:   &domains.EmailMetadata{NewValue: email1},
		Status:  &domains.StatusMetadata{NewValue: status1},
	}

	updateNewMetauserValues(metauser, tempMetauser)

	assert.Equal(t, metauser, expectedMetauser)
}

func TestMetaUsersImpl_isEqual_True(t *testing.T) {
	address1 := &domains.AddressValueMetadata{City: &domains.CityMetadata{ID: 1, Name: "1", IbgeID: 1}, Neighborhood: "neighborhood1"}
	address2 := &domains.AddressValueMetadata{City: &domains.CityMetadata{ID: 1, Name: "1", IbgeID: 1}, Neighborhood: "neighborhood1"}
	phone1 := &domains.PhoneValueMetadata{Phone: "phone1", UpdatedAt: time.Now(), ClientID: "clientid2"}
	phone2 := &domains.PhoneValueMetadata{Phone: "phone1", UpdatedAt: time.Now().Add(10), ClientID: "clientid1"}
	email1 := &domains.EmailValueMetadata{Email: "email1", UpdatedAt: time.Now()}
	email2 := &domains.EmailValueMetadata{Email: "email1", UpdatedAt: time.Now().Add(10)}
	status1 := &domains.StatusValueMetadata{Status: "status1", ClientID: "clientid1"}
	status2 := &domains.StatusValueMetadata{Status: "status1", ClientID: "clientid2"}

	metauser := domains.MetaUser{
		Address: &domains.AddressMetadata{NewValue: address1},
		Email:   &domains.EmailMetadata{NewValue: email1},
		Phones:  &domains.PhonesMetadata{NewValue: phone1},
		Status:  &domains.StatusMetadata{NewValue: status1},
	}

	tempMetauser := domains.MetaUser{
		Address: &domains.AddressMetadata{NewValue: address2},
		Email:   &domains.EmailMetadata{NewValue: email2},
		Phones:  &domains.PhonesMetadata{NewValue: phone2},
		Status:  &domains.StatusMetadata{NewValue: status2},
	}

	resp1 := isEqual(*metauser.Address.NewValue, *tempMetauser.Address.NewValue)
	assert.True(t, resp1)

	resp2 := isEqual(*metauser.Email.NewValue, *tempMetauser.Email.NewValue)
	assert.True(t, resp2)

	resp3 := isEqual(*metauser.Phones.NewValue, *tempMetauser.Phones.NewValue)
	assert.True(t, resp3)

	resp4 := isEqual(*metauser.Status.NewValue, *tempMetauser.Status.NewValue)
	assert.True(t, resp4)
}

func TestMetaUsersImpl_isEqual_False(t *testing.T) {
	address1 := &domains.AddressValueMetadata{City: &domains.CityMetadata{ID: 1, Name: "1", IbgeID: 1}, Neighborhood: "neighborhood1"}
	address2 := &domains.AddressValueMetadata{City: &domains.CityMetadata{ID: 1, Name: "1", IbgeID: 1}, Neighborhood: "neighborhood2"}
	phone1 := &domains.PhoneValueMetadata{Phone: "phone1", UpdatedAt: time.Now(), ClientID: "clientid1"}
	phone2 := &domains.PhoneValueMetadata{Phone: "phone2", UpdatedAt: time.Now(), ClientID: "clientid1"}
	email1 := &domains.EmailValueMetadata{Email: "email1", UpdatedAt: time.Now()}
	email2 := &domains.EmailValueMetadata{Email: "email2", UpdatedAt: time.Now()}
	status1 := &domains.StatusValueMetadata{Status: "status1", ClientID: "clientid1"}
	status2 := &domains.StatusValueMetadata{Status: "status2", ClientID: "clientid1"}
	username1 := &domains.UsernameValueMetadata{Username: "username1"}
	username2 := &domains.UsernameValueMetadata{Username: "username2", UpdatedAt: time.Now().Add(10), ClientID: "clientid1"}

	metauser := domains.MetaUser{
		Address:  &domains.AddressMetadata{NewValue: address1},
		Email:    &domains.EmailMetadata{NewValue: email1},
		Phones:   &domains.PhonesMetadata{NewValue: phone1},
		Status:   &domains.StatusMetadata{NewValue: status1},
		Username: &domains.UsernameMetadata{NewValue: username1},
	}

	tempMetauser := domains.MetaUser{
		Address:  &domains.AddressMetadata{NewValue: address2},
		Email:    &domains.EmailMetadata{NewValue: email2},
		Phones:   &domains.PhonesMetadata{NewValue: phone2},
		Status:   &domains.StatusMetadata{NewValue: status2},
		Username: &domains.UsernameMetadata{NewValue: username2},
	}

	resp1 := isEqual(*metauser.Address.NewValue, *tempMetauser.Address.NewValue)
	assert.False(t, resp1)

	resp2 := isEqual(*metauser.Email.NewValue, *tempMetauser.Email.NewValue)
	assert.False(t, resp2)

	resp3 := isEqual(*metauser.Phones.NewValue, *tempMetauser.Phones.NewValue)
	assert.False(t, resp3)

	resp4 := isEqual(*metauser.Status.NewValue, *tempMetauser.Status.NewValue)
	assert.False(t, resp4)

	resp5 := isEqual(*metauser.Username.NewValue, *tempMetauser.Username.NewValue)
	assert.False(t, resp5)
}

func TestMetaUsersImpl_parseUserToMetaUser(t *testing.T) {
	updatedAt := time.Now()
	clientId := "clientId"

	//Mocking user
	user := &domains.User{
		Phone: &domains.Phone{
			Cellphone:            "cellphone",
			Phone:                "phone",
			DDDCellphone:         "DddCellPhone",
			MobilePhoneConfirmed: true,
		},
		Address: &domains.Address{
			ID:            "id",
			City:          domains.City{ID: 1, IBGECityId: 1, Name: "city"},
			State:         domains.State{ID: 1, IBGEStateId: 1, Name: "state"},
			Country:       domains.Country{ID: 1, Initials: "initials", Name: "country"},
			Neighborhood:  "Neighborhood",
			Number:        1,
			AddressType:   "AddressType",
			Address1:      "Address1",
			Address2:      "Address2",
			Zipcode:       "Zipcode",
			SimpleAddress: "SimpleAddress",
		},
		Email:     "email",
		Username:  "username",
		Status:    "status",
		UpdatedAt: &updatedAt,
	}

	//Mocking metauser
	address := &domains.AddressValueMetadata{
		City:         &domains.CityMetadata{ID: 1, Name: "city", IbgeID: 1},
		State:        &domains.StateMetadata{ID: 1, Name: "state", IbgeID: 1},
		Country:      &domains.CountryMetadata{ID: 1, Name: "country"},
		Neighborhood: "Neighborhood",
		Number:       1,
		Address1:     "Address1",
		Address2:     "Address2",
		AddressType:  "AddressType",
		UpdatedAt:    updatedAt,
		ClientID:     clientId,
	}
	phone := &domains.PhoneValueMetadata{
		Phone:                "phone",
		CellPhone:            "cellphone",
		DddCellPhone:         "DddCellPhone",
		MobilePhoneConfirmed: true,
		UpdatedAt:            updatedAt,
		ClientID:             clientId,
	}
	email := &domains.EmailValueMetadata{Email: "email", UpdatedAt: updatedAt, ClientID: clientId}
	status := &domains.StatusValueMetadata{Status: "status", UpdatedAt: updatedAt, ClientID: clientId}
	username := &domains.UsernameValueMetadata{Username: "username", UpdatedAt: updatedAt, ClientID: clientId}
	password := &domains.PasswordValueMetadata{UpdatedAt: updatedAt, ClientID: clientId}

	expectedMetaUser := domains.MetaUser{
		Address:   &domains.AddressMetadata{NewValue: address},
		Email:     &domains.EmailMetadata{NewValue: email},
		Phones:    &domains.PhonesMetadata{NewValue: phone},
		Status:    &domains.StatusMetadata{NewValue: status},
		Username:  &domains.UsernameMetadata{NewValue: username},
		Password:  &domains.PasswordMetadata{NewValue: password},
		UpdatedAt: updatedAt,
	}

	metauser := parseUserToMetauser(user, clientId)
	assert.ObjectsAreEqual(metauser, expectedMetaUser)
}

func TestMetaUsersImpl_validateUpdatedAt_WhenItExist(t *testing.T) {
	updatedAt := time.Now()
	user := &domains.User{UpdatedAt: &updatedAt}

	validateUpdatedAt(user)

	assert.Equal(t, user.UpdatedAt, &updatedAt)
}

func TestMetaUsersImpl_validateUpdatedAt_WhenItDoesntExist(t *testing.T) {
	updatedAt := time.Now()
	user := &domains.User{}

	validateUpdatedAt(user)

	assert.NotNil(t, user.UpdatedAt)
	assert.NotEqual(t, user.UpdatedAt, &updatedAt)
}
