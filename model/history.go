package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type HistoryAction string

const (
	HistoryActionGet    HistoryAction = "get"
	HistoryActionCreate HistoryAction = "create"
	HistoryActionModify HistoryAction = "modify"
	HistoryActionDelete HistoryAction = "delete"
	HistoryActionSend   HistoryAction = "send"
)

type History struct {
	Identifier     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserId         primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Action         HistoryAction      `json:"action" bson:"action"`
	TargetQueryId  primitive.ObjectID `json:"target_query_id,omitempty" bson:"target_query_id,omitempty"`
	TargetMemberId primitive.ObjectID `json:"target_member_id,omitempty" bson:"target_member_id,omitempty"`
	TargetUserId   primitive.ObjectID `json:"target_user_id,omitempty" bson:"target_user_id,omitempty"`
	EMailContent   string             `json:"e_mail_content,omitempty" bson:"e_mail_content,omitempty"`
	DateTime       time.Time          `json:"date_time" bson:"date_time"`
}
