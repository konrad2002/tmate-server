package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Member struct {
	Identifier primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	Data       map[string]interface{} `bson:"data" json:"data"`
	CreatorId  primitive.ObjectID     `json:"creator_id,omitempty" bson:"creator_id,omitempty"`
	CreatedAt  time.Time              `json:"created_at" bson:"created_at"`
	ModifiedAt time.Time              `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
}
