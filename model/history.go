package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Action string

const (
	Get    Action = "get"
	Create Action = "create"
	Modify Action = "modify"
	Delete Action = "delete"
)

type History struct {
	Identifier     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserId         primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Action         Action             `json:"action" bson:"action"`
	TargetQueryId  primitive.ObjectID `json:"target_query_id,omitempty" bson:"target_query_id,omitempty"`
	TargetMemberId primitive.ObjectID `json:"target_member_id,omitempty" bson:"target_member_id,omitempty"`
}
