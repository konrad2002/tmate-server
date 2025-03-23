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

type HistoryRepository struct {
	mongoDB    *mongo.Database
	collection *mongo.Collection
}

func NewHistoryRepository(mongoDB *mongo.Database) HistoryRepository {
	return HistoryRepository{
		mongoDB:    mongoDB,
		collection: mongoDB.Collection("history"),
	}
}

func (mr *HistoryRepository) GetHistoriesByBsonDocument(d interface{}) ([]model.History, error) {

	queryOptions := options.FindOptions{}
	queryOptions.SetSort(bson.D{{"position", 1}})

	return mr.GetHistoriesByBsonDocumentWithOptions(d, &queryOptions)
}

func (mr *HistoryRepository) GetHistoriesByBsonDocumentWithOptions(d interface{}, queryOptions *options.FindOptions) ([]model.History, error) {
	var histories []model.History

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := mr.collection.Find(ctx, d, queryOptions)
	if err != nil {
		return []model.History{}, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var history model.History
		err := cursor.Decode(&history)
		if err != nil {
			return []model.History{}, err
		}
		histories = append(histories, history)
	}

	if err := cursor.Err(); err != nil {
		return []model.History{}, err
	}

	return histories, nil
}

func (mr *HistoryRepository) GetHistoryByBsonDocument(d interface{}) (model.History, error) {

	queryOptions := options.FindOptions{}
	queryOptions.SetSort(bson.D{{"position", 1}})

	return mr.GetHistoryByBsonDocumentWithOptions(d, &queryOptions)
}

func (mr *HistoryRepository) GetHistoryByBsonDocumentWithOptions(d interface{}, queryOptions *options.FindOptions) (model.History, error) {
	histories, err := mr.GetHistoriesByBsonDocumentWithOptions(d, queryOptions)

	if err != nil {
		return model.History{}, err
	}

	if len(histories) > 0 {
		return histories[0], nil
	}

	return model.History{}, errors.New("no entry found")
}

func (mr *HistoryRepository) SaveHistory(history model.History) (model.History, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	history.DateTime = time.Now()

	r, err := mr.collection.InsertOne(ctx, history)
	if err != nil {
		return model.History{}, err
	}

	return mr.GetHistoryByBsonDocument(bson.D{{"_id", r.InsertedID.(primitive.ObjectID)}})
}
