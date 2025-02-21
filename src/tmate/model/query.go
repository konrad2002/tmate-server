package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Query struct {
	Identifier primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty" bson:"name,omitempty"`
}
