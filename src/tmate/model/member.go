package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Member struct {
	ID   primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Data map[string]interface{} `bson:"data" json:"data"`
}
