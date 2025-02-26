package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Query struct {
	Identifier  primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Filter      bson.D             `json:"filter,omitempty" bson:"filter"`
	Projection  bson.D             `json:"projection,omitempty" bson:"projection"`
	Sort        bson.D             `json:"sort,omitempty" bson:"sort"`
	OwnerUserId primitive.ObjectID `json:"owner_user_id,omitempty" bson:"owner_user_id,omitempty"`
	Public      bool               `json:"public" bson:"public"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	ModifiedAt  time.Time          `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
}
