package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Member struct {
	Identifier primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	Data       map[string]interface{} `bson:"data" json:"data"`
}
