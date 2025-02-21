package repository

import "go.mongodb.org/mongo-driver/mongo"

type MemberRepository struct {
	mongoDB *mongo.Database
}

func NewMemberRepository(mongoDB *mongo.Database) MemberRepository {
	return MemberRepository{
		mongoDB: mongoDB,
	}
}
