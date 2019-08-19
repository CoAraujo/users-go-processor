package domains

import "time"

type User struct {
	ID                 string                                `json:"id,omitempty"`
	Email              string                                `json:"email,omitempty"`
	Username           string                                `json:"username,omitempty"`
	Name               string                                `json:"fullName,omitempty"`
	Gender             string                                `json:"gender,omitempty"`
	Status             string                                `json:"status,omitempty"`
	BirthDate          string                                `json:"birthDate,omitempty"`
	UpdatedAt          *time.Time                            `json:"updatedAt,omitempty"`
	Phone              *Phone                                `json:"phones,omitempty"`
	Address            *Address                              `json:"address,omitempty"`
}

type Phone struct {
	Cellphone            string `json:"cellPhone,omitempty"`
	DDDCellphone         string `json:"dddCellPhone,omitempty"`
	Phone                string `json:"phone,omitempty"`
	DDDPhone             string `json:"dddPhone,omitempty"`
	MobilePhoneConfirmed bool   `json:"mobilePhoneConfirmed,omitempty"`
}

type Address struct {
	City          City    `json:"city,omitempty"`
	State         State   `json:"state,omitempty"`
	Country       Country `json:"country,omitempty"`
	ID            string  `json:"id,omitempty"`
	Neighborhood  string  `json:"neighborhood,omitempty"`
	Zipcode       string  `json:"zipCode,omitempty"`
	Address1      string  `json:"address1,omitempty"`
	Address2      string  `json:"address2,omitempty"`
	SimpleAddress string  `json:"simpleAddress,omitempty"`
	AddressType   string  `json:"addressType,omitempty"`
	Number        int64   `json:"number,omitempty"`
}

type City struct {
	ID         int64  `json:"id,omitempty"`
	IBGECityId int64  `json:"ibgeId,omitempty"`
	Name       string `json:"name,omitempty"`
}

type State struct {
	Capital      int    `json:"capital,omitempty"`
	ID           int64  `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Initials     string `json:"initials,omitempty"`
	Abbreviation string `json:"abbreviation,omitempty"`
	IBGEStateId  int64  `json:"ibgeId,omitempty"`
}

type Country struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Initials string `json:"initials,omitempty"`
}