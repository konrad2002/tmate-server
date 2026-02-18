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

type CourseRepository struct {
	mongoDB    *mongo.Database
	collection *mongo.Collection
}

func NewCourseRepository(mongoDB *mongo.Database) CourseRepository {
	return CourseRepository{
		mongoDB:    mongoDB,
		collection: mongoDB.Collection("course"),
	}
}

func (cr *CourseRepository) GetCoursesByBsonDocument(d interface{}) ([]model.Course, error) {
	queryOptions := options.FindOptions{}
	queryOptions.SetSort(bson.D{{"time", -1}})

	return cr.GetCoursesByBsonDocumentWithOptions(d, &queryOptions)
}

func (cr *CourseRepository) GetCoursesByBsonDocumentWithOptions(d interface{}, queryOptions *options.FindOptions) ([]model.Course, error) {
	var courses []model.Course

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := cr.collection.Find(ctx, d, queryOptions)
	if err != nil {
		return []model.Course{}, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var course model.Course
		err := cursor.Decode(&course)
		if err != nil {
			return []model.Course{}, err
		}
		courses = append(courses, course)
	}

	if err := cursor.Err(); err != nil {
		return []model.Course{}, err
	}

	return courses, nil
}

func (cr *CourseRepository) GetCourseByBsonDocument(d interface{}) (model.Course, error) {
	queryOptions := options.FindOptions{}
	return cr.GetCourseByBsonDocumentWithOptions(d, &queryOptions)
}

func (cr *CourseRepository) GetCourseByBsonDocumentWithOptions(d interface{}, queryOptions *options.FindOptions) (model.Course, error) {
	courses, err := cr.GetCoursesByBsonDocumentWithOptions(d, queryOptions)

	if err != nil {
		return model.Course{}, err
	}

	if len(courses) > 0 {
		return courses[0], nil
	}

	return model.Course{}, errors.New("no entry found")
}

func (cr *CourseRepository) SaveCourse(course model.Course) (model.Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	course.CreatedAt = time.Now()
	course.ModifiedAt = time.Now()

	r, err := cr.collection.InsertOne(ctx, course)
	if err != nil {
		return model.Course{}, err
	}

	return cr.GetCourseByBsonDocument(bson.D{{"_id", r.InsertedID.(primitive.ObjectID)}})
}

func (cr *CourseRepository) UpdateCourse(course model.Course) (model.Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	course.ModifiedAt = time.Now()

	_, err := cr.collection.UpdateOne(
		ctx,
		bson.D{{"_id", course.Identifier}},
		bson.D{{"$set", course}},
	)
	if err != nil {
		return model.Course{}, err
	}

	return cr.GetCourseByBsonDocument(bson.D{{"_id", course.Identifier}})
}

func (cr *CourseRepository) DeleteCourse(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := cr.collection.DeleteOne(ctx, bson.D{{"_id", id}})
	return err
}
