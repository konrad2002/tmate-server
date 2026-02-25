package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Course struct {
	Identifier  primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Location    string             `json:"location,omitempty" bson:"location,omitempty"`
	Time        string             `json:"time,omitempty" bson:"time,omitempty"`
	Day         string             `json:"day,omitempty" bson:"day,omitempty"`
	Price       int                `json:"price,omitempty" bson:"position,omitempty"`
	TotalSpots  int                `json:"total_spots,omitempty" bson:"total_spots,omitempty"`
	FreeSpots   int                `json:"free_spots,omitempty" bson:"free_spots,omitempty"`
	Style       string             `json:"style,omitempty" bson:"style,omitempty"`
	Level       string             `json:"level,omitempty" bson:"level,omitempty"`
	Age         string             `json:"age,omitempty" bson:"age,omitempty"`
	Information string             `json:"information,omitempty" bson:"information,omitempty"`
	BeginDate   time.Time          `json:"begin_date,omitempty" bson:"begin_date,omitempty"`
	EndDate     time.Time          `json:"end_date,omitempty" bson:"end_date,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	ModifiedAt  time.Time          `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
}

type CourseRegistration struct {
	Identifier   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CourseId     primitive.ObjectID `json:"course_id,omitempty" bson:"course_id,omitempty"`
	RegisteredAt time.Time          `json:"registered_at" bson:"registered_at"`
}
