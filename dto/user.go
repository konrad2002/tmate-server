package dto

import (
	"github.com/konrad2002/tmate-server/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type LoginDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserDto struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"` // initial password, needs to be overwritten
	Role      string `json:"role"`
}

type UserInfoDto struct {
	Identifier   primitive.ObjectID `json:"id"`
	Username     string             `json:"username"`
	Email        string             `json:"email"`
	FirstName    string             `json:"first_name"`
	LastName     string             `json:"last_name"`
	Role         string             `json:"role"`
	Active       bool               `json:"active"`
	TempPassword bool               `json:"temp_password"`
	Logins       int                `json:"logins"`
	CreatedAt    time.Time          `json:"created_at"`
	ModifiedAt   time.Time          `json:"modified_at,omitempty"`
}

func UserToUserInfoDto(user model.User) UserInfoDto {
	return UserInfoDto{
		Identifier:   user.Identifier,
		Username:     user.Username,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Role:         user.Role,
		Active:       user.Active,
		TempPassword: user.TempPassword,
		Logins:       user.Logins,
		CreatedAt:    user.CreatedAt,
		ModifiedAt:   user.ModifiedAt,
	}
}

func UsersToUserInfoDtos(users []model.User) []UserInfoDto {
	var dtos []UserInfoDto
	for _, user := range users {
		dtos = append(dtos, UserToUserInfoDto(user))
	}
	return dtos
}
