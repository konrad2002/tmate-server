package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/konrad2002/tmate-server/dto"
	"github.com/konrad2002/tmate-server/model"
	"github.com/konrad2002/tmate-server/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(ur repository.UserRepository) UserService {
	return UserService{
		userRepository: ur,
	}
}

func (us *UserService) GetAll() ([]dto.UserInfoDto, error) {
	users, err := us.userRepository.GetUsersByBsonDocument(bson.D{})
	return dto.UsersToUserInfoDtos(users), err
}

func (us *UserService) GetUserById(id primitive.ObjectID) (dto.UserInfoDto, error) {
	user, err := us.userRepository.GetUserByBsonDocument(bson.D{{"_id", id}})
	return dto.UserToUserInfoDto(user), err
}

func (us *UserService) GetUserByUsername(username string) (dto.UserInfoDto, error) {
	user, err := us.userRepository.GetUserByBsonDocument(bson.D{{"username", username}})
	return dto.UserToUserInfoDto(user), err
}

func (us *UserService) CreateUser(user dto.CreateUserDto) (*dto.UserInfoDto, error) {
	existing, err := us.GetUserByUsername(user.Username)
	if err != nil {
		if err.Error() != repository.NoUserFoundError {
			fmt.Println("failed to lookup username:", err)
			return nil, err
		}
	}

	if !existing.Identifier.IsZero() {
		return nil, errors.New("username is already taken")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("failed to hash password:", err)
		return nil, err
	}

	newUser := model.User{
		Username:     user.Username,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Role:         user.Role,
		Active:       true,
		TempPassword: true,
		Logins:       0,
		Password:     string(passwordHash),
		CreatedAt:    time.Time{},
		ModifiedAt:   time.Time{},
	}

	createdUser, err := us.userRepository.SaveUser(newUser)
	if err != nil {
		fmt.Println("failed to save user:", err)
		return nil, err
	}

	userInfo := dto.UserToUserInfoDto(createdUser)

	return &userInfo, nil
}

func (us *UserService) Login(login dto.LoginDto) (string, error) {
	user, err := us.userRepository.GetUserByBsonDocument(bson.D{{"username", login.Username}})
	if err != nil {
		fmt.Println("failed to lookup username:", err)
		return "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		fmt.Println("failed to check password:", err)
		return "", errors.New("invalid password")
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.Identifier,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("TMATE_AUTH_SECRET")))

	if err != nil {
		fmt.Println("failed to generate token:", err)
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func (us *UserService) RemoveUser(id primitive.ObjectID) (dto.UserInfoDto, error) {
	user, err := us.userRepository.RemoveUserById(id)
	return dto.UserToUserInfoDto(user), err
}
