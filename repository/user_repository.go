package repository

import (
	"context"
	"errors"
	"github.com/konrad2002/tmate-server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var NoUserFoundError = "no user found"

type UserRepository struct {
	mongoDB    *mongo.Database
	collection *mongo.Collection
}

func NewUserRepository(mongoDB *mongo.Database) UserRepository {
	return UserRepository{
		mongoDB:    mongoDB,
		collection: mongoDB.Collection("user"),
	}
}

func (ur *UserRepository) GetUsersByBsonDocument(d interface{}) ([]model.User, error) {

	queryOptions := options.FindOptions{}
	queryOptions.SetSort(bson.D{{"username", 1}})

	return ur.GetUsersByBsonDocumentWithOptions(d, &queryOptions)
}

func (ur *UserRepository) GetUsersByBsonDocumentWithOptions(d interface{}, queryOptions *options.FindOptions) ([]model.User, error) {
	var users []model.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := ur.collection.Find(ctx, d, queryOptions)
	if err != nil {
		return []model.User{}, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user model.User
		err := cursor.Decode(&user)
		if err != nil {
			return []model.User{}, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return []model.User{}, err
	}

	return users, nil
}

func (ur *UserRepository) GetUserByBsonDocument(d interface{}) (model.User, error) {

	queryOptions := options.FindOptions{}
	queryOptions.SetSort(bson.D{{"username", 1}})

	return ur.GetUserByBsonDocumentWithOptions(d, &queryOptions)
}

func (ur *UserRepository) GetUserByBsonDocumentWithOptions(d interface{}, queryOptions *options.FindOptions) (model.User, error) {
	users, err := ur.GetUsersByBsonDocumentWithOptions(d, queryOptions)

	if err != nil {
		return model.User{}, err
	}

	if len(users) > 0 {
		return users[0], nil
	}

	return model.User{}, errors.New(NoUserFoundError)
}

func (ur *UserRepository) SaveUser(user model.User) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user.CreatedAt = time.Now()
	user.ModifiedAt = time.Now()

	r, err := ur.collection.InsertOne(ctx, user)
	if err != nil {
		return model.User{}, err
	}

	return ur.GetUserByBsonDocument(bson.D{{"_id", r.InsertedID.(primitive.ObjectID)}})
}

func (ur *UserRepository) UpdateUser(user model.User) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user.ModifiedAt = time.Now()

	_, err := ur.collection.ReplaceOne(ctx, bson.D{{"_id", user.Identifier}}, user)
	if err != nil {
		return model.User{}, err
	}

	return ur.GetUserByBsonDocument(bson.D{{"_id", user.Identifier}})
}

func (ur *UserRepository) RemoveUserById(id primitive.ObjectID) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := ur.GetUserByBsonDocument(bson.D{{"_id", id}})
	if err != nil {
		return model.User{}, err
	}

	_, err = ur.collection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
