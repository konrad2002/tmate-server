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
	Role         string             `json:"role,omitempty" bson:"role"`
	Active       bool               `json:"active,omitempty" bson:"active"`
	TempPassword bool               `json:"temp_password,omitempty" bson:"temp_password"`
	Logins       int                `json:"logins,omitempty" bson:"logins"`
	Token        string             `json:"token,omitempty" bson:"token,omitempty"`
	Password     string             `json:"password,omitempty" bson:"password,omitempty"`
	Permissions  Permission         `json:"permissions,omitempty" bson:"permissions,omitempty"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	ModifiedAt   time.Time          `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
	LastLoginAt  time.Time          `json:"last_login_at,omitempty" bson:"last_login_at,omitempty"`
}

type Role struct {
	Name        string     `json:"name" bson:"name"`
	Permissions Permission `json:"permissions" bson:"permissions,omitempty"`
}

type PermissionLevel int

const (
	PermissionLevelNone   PermissionLevel = iota
	PermissionLevelRead   PermissionLevel = iota
	PermissionLevelWrite  PermissionLevel = iota
	PermissionLevelDelete PermissionLevel = iota
)

type Permission struct {
	// int fields: 0; 1=read; 2=write; 3=delete
	SuperUser                bool                      `json:"super_user" bson:"super_user"`                                 // used only for config init, does not overwrite other perms
	UserManagement           bool                      `json:"user_management" bson:"user_management"`                       // used
	TableStructureManagement bool                      `json:"table_structure_management" bson:"table_structure_management"` // unused
	EmailAddressManagement   bool                      `json:"email_address_management" bson:"email_address_management"`     // unused
	EmailAddressUsage        map[string]bool           `json:"email_address_usage" bson:"email_address_usage"`               // used
	BypassEmailRegex         bool                      `json:"bypass_email_regex" bson:"bypass_email_regex"`                 // unused
	QueryManagement          bool                      `json:"query_management" bson:"query_management"`                     // used
	MemberAdmin              PermissionLevel           `json:"member_admin" bson:"member_admin"`                             // used
	Member                   map[string]map[string]int `json:"member" bson:"member"`                                         // unused
	// 						 group x column
}
