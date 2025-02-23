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

type MemberRepository struct {
	mongoDB    *mongo.Database
	collection *mongo.Collection
}

func NewMemberRepository(mongoDB *mongo.Database) MemberRepository {
	return MemberRepository{
		mongoDB:    mongoDB,
		collection: mongoDB.Collection("member"),
	}
}

func (mr *MemberRepository) GetMembersByBsonDocument(d interface{}) ([]model.Member, error) {

	queryOptions := options.FindOptions{}
	queryOptions.SetSort(bson.D{{"position", 1}})

	return mr.GetMembersByBsonDocumentWithOptions(d, &queryOptions)
}

func (mr *MemberRepository) GetMembersByBsonDocumentWithOptions(d interface{}, queryOptions *options.FindOptions) ([]model.Member, error) {
	var members []model.Member

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := mr.collection.Find(ctx, d, queryOptions)
	if err != nil {
		return []model.Member{}, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var member model.Member
		err := cursor.Decode(&member)
		if err != nil {
			return []model.Member{}, err
		}
		members = append(members, member)
	}

	if err := cursor.Err(); err != nil {
		return []model.Member{}, err
	}

	return members, nil
}

func (mr *MemberRepository) GetMemberByBsonDocument(d interface{}) (model.Member, error) {

	queryOptions := options.FindOptions{}
	queryOptions.SetSort(bson.D{{"position", 1}})

	return mr.GetMemberByBsonDocumentWithOptions(d, &queryOptions)
}

func (mr *MemberRepository) GetMemberByBsonDocumentWithOptions(d interface{}, queryOptions *options.FindOptions) (model.Member, error) {
	members, err := mr.GetMembersByBsonDocumentWithOptions(d, queryOptions)

	if err != nil {
		return model.Member{}, err
	}

	if len(members) > 0 {
		return members[0], nil
	}

	return model.Member{}, errors.New("no entry found")
}

func (mr *MemberRepository) SaveMember(member model.Member) (model.Member, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	member.CreatedAt = time.Now()
	member.ModifiedAt = time.Now()

	r, err := mr.collection.InsertOne(ctx, member)
	if err != nil {
		return model.Member{}, err
	}

	return mr.GetMemberByBsonDocument(bson.D{{"_id", r.InsertedID.(primitive.ObjectID)}})
}

func (mr *MemberRepository) UpdateMember(member model.Member) (model.Member, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	member.ModifiedAt = time.Now()

	_, err := mr.collection.ReplaceOne(ctx, bson.D{{"_id", member.Identifier}}, member)
	if err != nil {
		return model.Member{}, err
	}

	return mr.GetMemberByBsonDocument(bson.D{{"_id", member.Identifier}})
}
