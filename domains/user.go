package domains

import "time"

type User struct {
	ID        string    `bson:"_id,omitempty" json:"_id,omitempty"`
	Email     string    `bson:"email,omitempty" json:"email,omitempty"`
	Username  string    `bson:"username,omitempty" json:"username,omitempty"`
	Name      string    `bson:"fullName,omitempty" json:"fullName,omitempty"`
	Gender    string    `bson:"gender,omitempty" json:"gender,omitempty"`
	Status    string    `bson:"status,omitempty" json:"status,omitempty"`
	BirthDate string    `bson:"birthDate,omitempty" json:"birthDate,omitempty"`
	Phones    *Phone    `bson:"phones,omitempty" json:"phones,omitempty"`
	ClientID  string    `bson:"clientId,omitempty" json:"clientId,omitempty"`
	UpdatedAt time.Time `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type Phone struct {
	Phone                string    `bson:"phone,omitempty" json:"phone,omitempty"`
	CellPhone            string    `bson:"cellphone,omitempty" json:"cellphone,omitempty"`
	DddCellPhone         string    `bson:"ddd_cellphone,omitempty" json:"ddd_cellphone,omitempty"`
	MobilePhoneConfirmed bool      `bson:"mobile_phone_confirmed,omitempty" json:"mobile_phone_confirmed,omitempty"`
	UpdatedAt            time.Time `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
