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

type QueryRepository struct {
	mongoDB    *mongo.Database
	collection *mongo.Collection
}

func NewQueryRepository(mongoDB *mongo.Database) QueryRepository {
	return QueryRepository{
		mongoDB:    mongoDB,
		collection: mongoDB.Collection("query"),
	}
}

func (qr *QueryRepository) GetQueriesByBsonDocument(d interface{}) ([]model.Query, error) {

	queryOptions := options.FindOptions{}
	queryOptions.SetSort(bson.D{})

	return qr.GetQueriesByBsonDocumentWithOptions(d, &queryOptions)
}

func (qr *QueryRepository) GetQueriesByBsonDocumentWithOptions(d interface{}, queryOptions *options.FindOptions) ([]model.Query, error) {
	var queries []model.Query

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := qr.collection.Find(ctx, d, queryOptions)
	if err != nil {
		return []model.Query{}, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var query model.Query
		err := cursor.Decode(&query)
		if err != nil {
			return []model.Query{}, err
		}
		queries = append(queries, query)
	}

	if err := cursor.Err(); err != nil {
		return []model.Query{}, err
	}

	return queries, nil
}

func (qr *QueryRepository) GetQueryByBsonDocument(d interface{}) (model.Query, error) {

	queryOptions := options.FindOptions{}
	queryOptions.SetSort(bson.D{{"position", 1}})

	return qr.GetQueryByBsonDocumentWithOptions(d, &queryOptions)
}

func (qr *QueryRepository) GetQueryByBsonDocumentWithOptions(d interface{}, queryOptions *options.FindOptions) (model.Query, error) {
	queries, err := qr.GetQueriesByBsonDocumentWithOptions(d, queryOptions)

	if err != nil {
		return model.Query{}, err
	}

	if len(queries) > 0 {
		return queries[0], nil
	}

	return model.Query{}, errors.New("no entry found")
}

func (qr *QueryRepository) Save(query model.Query) (model.Query, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query.ModifiedAt = time.Now()

	r, err := qr.collection.InsertOne(ctx, query)
	if err != nil {
		return model.Query{}, err
	}

	return qr.GetQueryByBsonDocument(bson.D{{"_id", r.InsertedID.(primitive.ObjectID)}})
}
