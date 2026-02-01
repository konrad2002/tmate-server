package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type FieldType string

const (
	String      FieldType = "string"
	Number      FieldType = "number"
	Email       FieldType = "email"
	Select      FieldType = "select"
	PhoneNumber FieldType = "phone_number"
	MultiSelect FieldType = "multi_select"
	Boolean     FieldType = "boolean"
	Date        FieldType = "date"
	Family      FieldType = "family"
)

func GetAllFieldType() []FieldType {
	return []FieldType{
		String,
		Number,
		Email,
		Select,
		PhoneNumber,
		MultiSelect,
		Boolean,
		Date,
		Family,
	}
}

type Field struct {
	Identifier  primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	DisplayName string             `json:"display_name,omitempty" bson:"display_name,omitempty"`
	Type        FieldType          `json:"type,omitempty" bson:"type,omitempty"`
	Data        *FieldData         `json:"data,omitempty" bson:"data,omitempty"`
	Nullable    bool               `json:"nullable" bson:"nullable"`
	Position    int                `json:"position" bson:"position"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	ModifiedAt  time.Time          `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
}

type FieldData struct {
	Options   map[string]string `json:"options,omitempty" bson:"options,omitempty"`
	Validator *string           `json:"validator,omitempty" bson:"validator,omitempty"`
}
