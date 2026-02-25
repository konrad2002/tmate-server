package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SpecialForm string

const (
	DefaultForm SpecialForm = "default"
	CourseForm  SpecialForm = "course"
)

type Form struct {
	Identifier  primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	FormGroups  []FormGroup        `json:"form_groups" bson:"form_groups"`
	Defaults    []FormDefault      `json:"defaults" bson:"defaults"`
	SpecialForm SpecialForm        `json:"special_form" bson:"special_form"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	ModifiedAt  time.Time          `json:"modified_at,omitempty" bson:"modified_at,omitempty"`
}

type FormGroup struct {
	Name   string   `json:"name" bson:"name"`
	Fields []string `json:"fields" bson:"fields"`
}

type FormDefault struct {
	Field string `json:"field" bson:"field"`
	Value string `json:"value" bson:"value"`
}
