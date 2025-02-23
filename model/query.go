package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Query struct {
	Identifier  primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Query       bson.D             `json:"query,omitempty" bson:"query,omitempty"`
	OwnerUserId primitive.ObjectID `json:"owner_user_id,omitempty" bson:"owner_user_id,omitempty"`
	Public      bool               `json:"public" bson:"public"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	ModifiedAt  time.Time          `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
}
