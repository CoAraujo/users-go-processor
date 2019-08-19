package metauser

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
	metausersCollection = "metausers"
)

var (
	instance MetaUsers
	once     sync.Once
)

type MetaUsers interface {
	Get(id string) (*domains.MetaUser, error)
	Insert(user *domains.User, clientID string) (string, error)
	Update(newUser *domains.User, oldUser *domains.MetaUser, clientID string) error
	Delete(id string) error
}

type metaUsersImpl struct{}

func GetInstance() MetaUsers {
	once.Do(func() {
		instance = &metaUsersImpl{}
	})
	return instance
}

func (m *metaUsersImpl) Get(id string) (*domains.MetaUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var metauser domains.MetaUser
	if mgoErr := storage.GetInstance().FindOne(ctx, metausersCollection, map[string]interface{}{"_id": id}, &metauser); mgoErr != nil {
		return nil, mgoErr
	}

	return &metauser, nil
}

func (m *metaUsersImpl) Insert(user *domains.User, clientID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	metauser := parseUserToMetauser(user, clientID)

	id, mgoErr := storage.GetInstance().Insert(ctx, metausersCollection, &metauser)
	if mgoErr != nil {
		return id.(string), mgoErr
	}

	return id.(string), nil
}

func (m *metaUsersImpl) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if mgoErr := storage.GetInstance().Remove(ctx, metausersCollection, map[string]interface{}{"_id": id}); mgoErr != nil {
		return mgoErr
	}

	return nil
}

func (m *metaUsersImpl) Update(user *domains.User, metauser *domains.MetaUser, clientID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	tempMetauser := parseUserToMetauser(user, clientID)
	updateNewMetauserValues(metauser, tempMetauser)

	if mgoErr := storage.GetInstance().Update(ctx, metausersCollection, map[string]interface{}{"_id": metauser.ID},
		map[string]interface{}{"$set": &metauser}); mgoErr != nil {
		return mgoErr
	}

	return nil
}

var updateNewMetauserValues = func(metauser *domains.MetaUser, tempMetauser domains.MetaUser) {
	if metauser.Address != nil && tempMetauser.Address != nil && !isEqual(*metauser.Address.NewValue, *tempMetauser.Address.NewValue) {
		metauser.Address.LastValue = metauser.Address.NewValue
		metauser.Address.NewValue = tempMetauser.Address.NewValue
	} else if metauser.Address == nil && tempMetauser.Address != nil {
		metauser.Address = tempMetauser.Address
	}

	if metauser.Phones != nil && tempMetauser.Phones != nil && !isEqual(*metauser.Phones.NewValue, *tempMetauser.Phones.NewValue) {
		metauser.Phones.LastValue = metauser.Phones.NewValue
		metauser.Phones.NewValue = tempMetauser.Phones.NewValue
	} else if metauser.Phones == nil && tempMetauser.Phones != nil {
		metauser.Phones = tempMetauser.Phones
	}

	if metauser.Email != nil && tempMetauser.Email != nil && !isEqual(*metauser.Email.NewValue, *tempMetauser.Email.NewValue) {
		metauser.Email.LastValue = metauser.Email.NewValue
		metauser.Email.NewValue = tempMetauser.Email.NewValue
	} else if metauser.Email == nil && tempMetauser.Email != nil {
		metauser.Email = tempMetauser.Email
	}

	if metauser.Username != nil && tempMetauser.Username != nil && !isEqual(*metauser.Username.NewValue, *tempMetauser.Username.NewValue) {
		metauser.Username.LastValue = metauser.Username.NewValue
		metauser.Username.NewValue = tempMetauser.Username.NewValue
	} else if metauser.Username == nil && tempMetauser.Username != nil {
		metauser.Username = tempMetauser.Username
	}

	if metauser.Status != nil && tempMetauser.Status != nil && !isEqual(*metauser.Status.NewValue, *tempMetauser.Status.NewValue) {
		metauser.Status.LastValue = metauser.Status.NewValue
		metauser.Status.NewValue = tempMetauser.Status.NewValue
	} else if metauser.Status == nil && tempMetauser.Status != nil {
		metauser.Status = tempMetauser.Status
	}
}

var isEqual = func(interface1, interface2 interface{}) bool {
	return cmp.Equal(interface1, interface2, cmpopts.IgnoreFields(interface1, "UpdatedAt", "ClientID"))
}

var parseUserToMetauser = func(user *domains.User, clientID string) domains.MetaUser {
	validateUpdatedAt(user)

	var phone = &domains.PhoneValueMetadata{}
	if user.Phone != nil {
		phone = &domains.PhoneValueMetadata{
			Phone:                user.Phone.Phone,
			CellPhone:            user.Phone.Cellphone,
			DddCellPhone:         user.Phone.DDDCellphone,
			MobilePhoneConfirmed: user.Phone.MobilePhoneConfirmed,
			UpdatedAt:            *user.UpdatedAt,
			ClientID:             clientID,
		}
	}

	var address = &domains.AddressValueMetadata{}
	if user.Address != nil {
		address = &domains.AddressValueMetadata{
			City: &domains.CityMetadata{
				ID:     user.Address.City.ID,
				IbgeID: user.Address.City.IBGECityId,
				Name:   user.Address.City.Name,
			},
			State: &domains.StateMetadata{
				ID:           user.Address.State.ID,
				IbgeID:       user.Address.State.IBGEStateId,
				Abbreviation: user.Address.State.Abbreviation,
				Name:         user.Address.State.Name,
			},
			Country: &domains.CountryMetadata{
				ID:   user.Address.Country.ID,
				Name: user.Address.Country.Name,
			},
			Neighborhood: user.Address.Neighborhood,
			ZipCode:      user.Address.Zipcode,
			Address1:     user.Address.Address1,
			Address2:     user.Address.Address2,
			AddressType:  user.Address.AddressType,
			Number:       user.Address.Number,
			UpdatedAt:    *user.UpdatedAt,
			ClientID:     clientID,
		}
	}

	metauser := domains.MetaUser{
		ID: user.ID,
		Phones: &domains.PhonesMetadata{
			NewValue: phone,
		},
		Address: &domains.AddressMetadata{
			NewValue: address,
		},
		Password: &domains.PasswordMetadata{
			NewValue: &domains.PasswordValueMetadata{
				UpdatedAt: *user.UpdatedAt,
				ClientID:  clientID,
			},
		},
		Email: &domains.EmailMetadata{
			NewValue: &domains.EmailValueMetadata{
				Email:     user.Email,
				UpdatedAt: *user.UpdatedAt,
				ClientID:  clientID,
			},
		},
		Username: &domains.UsernameMetadata{
			NewValue: &domains.UsernameValueMetadata{
				Username:  user.Username,
				UpdatedAt: *user.UpdatedAt,
				ClientID:  clientID,
			},
		},
		Status: &domains.StatusMetadata{
			NewValue: &domains.StatusValueMetadata{
				Status:    user.Status,
				UpdatedAt: *user.UpdatedAt,
				ClientID:  clientID,
			},
		},
		UpdatedAt: *user.UpdatedAt,
	}
	return metauser
}

var validateUpdatedAt = func(user *domains.User) {
	if user.UpdatedAt == nil {
		log.Errorf("[Meta_User validateUpdatedAt] user without updatedAt. Setting updatedAt with time.Now(). Id: %s ", user.ID)
		updatedAt := time.Now()
		user.UpdatedAt = &updatedAt
	}
}
