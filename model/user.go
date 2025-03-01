package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Identifier   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username     string             `json:"username,omitempty" bson:"username,omitempty"`
	FirstName    string             `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName     string             `json:"last_name,omitempty" bson:"last_name,omitempty"`
	MemberId     primitive.ObjectID `json:"member_id" bson:"member_id,omitempty"`
	Member       Member             `json:"member,omitempty" bson:"-"`
	Role         string             `json:"role,omitempty" bson:"role,omitempty"`
	Active       bool               `json:"active,omitempty" bson:"active,omitempty"`
	TempPassword bool               `json:"temp_password,omitempty" bson:"temp_password,omitempty"`
	Logins       int                `json:"logins,omitempty" bson:"logins,omitempty"`
	Token        string             `json:"token,omitempty" bson:"token,omitempty"`
	Password     string             `json:"password,omitempty" bson:"password,omitempty"`
}
