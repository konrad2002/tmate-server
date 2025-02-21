package repository

import (
	"context"
	"errors"
	"github.com/konrad2002/tmate-server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type FieldRepository struct {
	mongoDB    *mongo.Database
	collection *mongo.Collection
}

func NewFieldRepository(mongoDB *mongo.Database) FieldRepository {
	return FieldRepository{
		mongoDB:    mongoDB,
		collection: mongoDB.Collection("field"),
	}
}

func (fr *FieldRepository) GetFieldsByBsonDocument(d interface{}) ([]model.Field, error) {

	queryOptions := options.FindOptions{}
	queryOptions.SetSort(bson.D{{"position", 1}})

	return fr.GetFieldsByBsonDocumentWithOptions(d, &queryOptions)
}

func (fr *FieldRepository) GetFieldsByBsonDocumentWithOptions(d interface{}, queryOptions *options.FindOptions) ([]model.Field, error) {
	var fields []model.Field

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := fr.collection.Find(ctx, d, queryOptions)
	if err != nil {
		return []model.Field{}, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var field model.Field
		err := cursor.Decode(&field)
		if err != nil {
			return []model.Field{}, err
		}
		fields = append(fields, field)
	}

	if err := cursor.Err(); err != nil {
		return []model.Field{}, err
	}

	return fields, nil
}

func (fr *FieldRepository) GetFieldByBsonDocument(d interface{}) (model.Field, error) {

	queryOptions := options.FindOptions{}
	queryOptions.SetSort(bson.D{{"position", 1}})

	return fr.GetFieldByBsonDocumentWithOptions(d, &queryOptions)
}

func (fr *FieldRepository) GetFieldByBsonDocumentWithOptions(d interface{}, queryOptions *options.FindOptions) (model.Field, error) {
	fields, err := fr.GetFieldsByBsonDocumentWithOptions(d, queryOptions)

	if err != nil {
		return model.Field{}, err
	}

	if len(fields) > 0 {
		return fields[0], nil
	}

	return model.Field{}, errors.New("no entry found")
}
