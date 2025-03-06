package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Identifier   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username     string             `json:"username" bson:"username"`
	Email        string             `json:"email" bson:"email"`
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
	Permissions  Permission         `json:"permissions,omitempty" bson:"permissions,omitempty"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	ModifiedAt   time.Time          `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
}

type Role struct {
	Name        string     `json:"name" bson:"name"`
	Permissions Permission `json:"permissions" bson:"permissions,omitempty"`
}

type Permission struct {
	// int fields: 0; 1=read; 2=write; 3=delete
	UserManagement           bool                      `json:"user_management" bson:"user_management"`
	TableStructureManagement bool                      `json:"table_structure_management" bson:"table_structure_management"`
	EmailAddressManagement   bool                      `json:"email_address_management" bson:"email_address_management"`
	EmailAddressUsage        map[string]bool           `json:"email_address_usage" bson:"email_address_usage"`
	ByPassEmailRegex         bool                      `json:"by_pass_email_regex" bson:"by_pass_email_regex"`
	QueryManagement          bool                      `json:"query_management" bson:"query_management"`
	MemberAdmin              int                       `json:"member_admin" bson:"member_admin"`
	Member                   map[string]map[string]int `json:"member" bson:"member"`
	// 						 group x column
}
