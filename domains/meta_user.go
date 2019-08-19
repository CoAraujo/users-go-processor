package domains

import "time"

type MetaUser struct {
	ID        string            `bson:"_id" json:"_id"`
	Phones    *PhonesMetadata   `bson:"phones,omitempty" json:"phones,omitempty"`
	Address   *AddressMetadata  `bson:"address,omitempty" json:"address,omitempty"`
	Password  *PasswordMetadata `bson:"password,omitempty" json:"password,omitempty"`
	Email     *EmailMetadata    `bson:"email,omitempty" json:"email,omitempty"`
	Username  *UsernameMetadata `bson:"username,omitempty" json:"username,omitempty"`
	Status    *StatusMetadata   `bson:"status,omitempty" json:"status,omitempty"`
	UpdatedAt time.Time         `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type PhonesMetadata struct {
	LastValue *PhoneValueMetadata `bson:"last_value,omitempty" json:"last_value,omitempty"`
	NewValue  *PhoneValueMetadata `bson:"new_value" json:"new_value"`
}

type PhoneValueMetadata struct {
	Phone                string    `bson:"phone,omitempty" json:"phone,omitempty"`
	CellPhone            string    `bson:"cellphone,omitempty" json:"cellphone,omitempty"`
	DddCellPhone         string    `bson:"ddd_cellphone,omitempty" json:"ddd_cellphone,omitempty"`
	MobilePhoneConfirmed bool      `bson:"mobile_phone_confirmed,omitempty" json:"mobile_phone_confirmed,omitempty"`
	ClientID             string    `bson:"client_id,omitempty" json:"client_id,omitempty"`
	UpdatedAt            time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type AddressMetadata struct {
	LastValue *AddressValueMetadata `bson:"last_value,omitempty" json:"last_value,omitempty"`
	NewValue  *AddressValueMetadata `bson:"new_value" json:"new_value"`
}

type AddressValueMetadata struct {
	City         *CityMetadata    `bson:"city,omitempty" json:"city,omitempty"`
	State        *StateMetadata   `bson:"state,omitempty" json:"state,omitempty"`
	Country      *CountryMetadata `bson:"country,omitempty" json:"country,omitempty"`
	Neighborhood string           `bson:"neighborhood,omitempty" json:"neighborhood,omitempty"`
	ZipCode      string           `bson:"zipCode,omitempty" json:"zipCode,omitempty"`
	Address1     string           `bson:"address1,omitempty" json:"address1,omitempty"`
	Address2     string           `bson:"address2,omitempty" json:"address2,omitempty"`
	AddressType  string           `bson:"addressType,omitempty" json:"addressType,omitempty"`
	Number       int64            `bson:"number,omitempty" json:"number,omitempty"`
	UpdatedAt    time.Time        `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	ClientID     string           `bson:"client_id,omitempty" json:"client_id,omitempty"`
}

type CityMetadata struct {
	ID     int64  `bson:"id,omitempty" json:"id,omitempty"`
	IbgeID int64  `bson:"ibge_id,omitempty" json:"ibge_id,omitempty"`
	Name   string `bson:"name,omitempty" json:"name,omitempty"`
}

type StateMetadata struct {
	ID           int64  `bson:"id,omitempty" json:"id,omitempty"`
	IbgeID       int64  `bson:"ibge_id,omitempty" json:"ibge_id,omitempty"`
	Name         string `bson:"name,omitempty" json:"name,omitempty"`
	Abbreviation string `bson:"abbreviation,omitempty" json:"abbreviation,omitempty"`
}

type CountryMetadata struct {
	ID   int64  `bson:"id,omitempty" json:"id,omitempty"`
	Name string `bson:"name,omitempty" json:"name,omitempty"`
}

type PasswordMetadata struct {
	LastValue *PasswordValueMetadata `bson:"last_value,omitempty" json:"last_value,omitempty"`
	NewValue  *PasswordValueMetadata `bson:"new_value" json:"new_value"`
}

type PasswordValueMetadata struct {
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	ClientID  string    `bson:"client_id,omitempty" json:"client_id,omitempty"`
}

type EmailMetadata struct {
	LastValue *EmailValueMetadata `bson:"last_value,omitempty" json:"last_value,omitempty"`
	NewValue  *EmailValueMetadata `bson:"new_value" json:"new_value"`
}

type EmailValueMetadata struct {
	Email     string    `bson:"email,omitempty" json:"email,omitempty"`
	ClientID  string    `bson:"client_id,omitempty" json:"client_id,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type UsernameMetadata struct {
	LastValue *UsernameValueMetadata `bson:"last_value,omitempty" json:"last_value,omitempty"`
	NewValue  *UsernameValueMetadata `bson:"new_value" json:"new_value"`
}

type UsernameValueMetadata struct {
	Username  string    `bson:"username,omitempty" json:"username,omitempty"`
	ClientID  string    `bson:"client_id,omitempty" json:"client_id,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type StatusMetadata struct {
	LastValue *StatusValueMetadata `bson:"last_value,omitempty" json:"last_value,omitempty"`
	NewValue  *StatusValueMetadata `bson:"new_value" json:"new_value"`
}

type StatusValueMetadata struct {
	Status    string    `bson:"status,omitempty" json:"status,omitempty"`
	ClientID  string    `bson:"client_id,omitempty" json:"client_id,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
