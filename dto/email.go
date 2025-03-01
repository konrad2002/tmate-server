package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EmailSenderDto struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type Attachment struct {
	Filename string `json:"filename"`
	Content  []byte `json:"content"`
}

type SendEmailDto struct {
	Sender       string               `json:"sender"`
	Receivers    []primitive.ObjectID `json:"receivers"`
	Subject      string               `json:"subject"`
	BodyTemplate string               `json:"body_template"`
	Attachments  []Attachment         `json:"attachments"`
}
