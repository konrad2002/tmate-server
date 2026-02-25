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

type FormRepository struct {
	mongoDB    *mongo.Database
	collection *mongo.Collection
}

func NewFormRepository(mongoDB *mongo.Database) FormRepository {
	return FormRepository{
		mongoDB:    mongoDB,
		collection: mongoDB.Collection("form"),
	}
}

func (fr *FormRepository) GetFormsByBsonDocument(d interface{}) ([]model.Form, error) {

	formOptions := options.FindOptions{}
	formOptions.SetSort(bson.D{})

	return fr.GetFormsByBsonDocumentWithOptions(d, &formOptions)
}

func (fr *FormRepository) GetFormsByBsonDocumentWithOptions(d interface{}, formOptions *options.FindOptions) ([]model.Form, error) {
	var forms []model.Form

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := fr.collection.Find(ctx, d, formOptions)
	if err != nil {
		return []model.Form{}, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var form model.Form
		err := cursor.Decode(&form)
		if err != nil {
			return []model.Form{}, err
		}
		forms = append(forms, form)
	}

	if err := cursor.Err(); err != nil {
		return []model.Form{}, err
	}

	return forms, nil
}

func (fr *FormRepository) GetFormByBsonDocument(d interface{}) (model.Form, error) {

	formOptions := options.FindOptions{}
	formOptions.SetSort(bson.D{{"position", 1}})

	return fr.GetFormByBsonDocumentWithOptions(d, &formOptions)
}

func (fr *FormRepository) GetFormByBsonDocumentWithOptions(d interface{}, formOptions *options.FindOptions) (model.Form, error) {
	forms, err := fr.GetFormsByBsonDocumentWithOptions(d, formOptions)

	if err != nil {
		return model.Form{}, err
	}

	if len(forms) > 0 {
		return forms[0], nil
	}

	return model.Form{}, errors.New("no entry found")
}

func (fr *FormRepository) SaveForm(form model.Form) (model.Form, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	form.CreatedAt = time.Now()
	form.ModifiedAt = time.Now()

	r, err := fr.collection.InsertOne(ctx, form)
	if err != nil {
		return model.Form{}, err
	}

	return fr.GetFormByBsonDocument(bson.D{{"_id", r.InsertedID.(primitive.ObjectID)}})
}

func (fr *FormRepository) UpdateForm(form model.Form) (model.Form, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	form.ModifiedAt = time.Now()

	_, err := fr.collection.ReplaceOne(ctx, bson.D{{"_id", form.Identifier}}, form)
	if err != nil {
		return model.Form{}, err
	}

	return fr.GetFormByBsonDocument(bson.D{{"_id", form.Identifier}})
}

func (fr *FormRepository) RemoveFormById(id primitive.ObjectID) (model.Form, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	form, err := fr.GetFormByBsonDocument(bson.D{{"_id", id}})
	if err != nil {
		return model.Form{}, err
	}

	_, err = fr.collection.DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return model.Form{}, err
	}
	return form, nil
}
