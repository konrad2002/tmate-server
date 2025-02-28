package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type SpecialFields struct {
	FirstName      string        `json:"first_name"`
	LastName       string        `json:"last_name"`
	EMail          string        `json:"e_mail"`
	EMail2         string        `json:"e_mail_2"`
	Family         string        `json:"family"`
	AttestDate     string        `json:"attest_date"`
	AttestRequired string        `json:"attest_required"`
	Address        AddressFields `json:"address"`
}

type AddressFields struct {
	Street     string `json:"street"`
	Number     string `json:"number"`
	City       string `json:"city"`
	PostalCode string `json:"postal_code"`
}

type Config struct {
	DefaultQuery primitive.ObjectID `json:"default_query"`
}
