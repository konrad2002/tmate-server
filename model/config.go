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
	Courses        string        `json:"courses"`
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

type EmailConfig struct {
	Address string          `json:"address"`
	Name    string          `json:"name"`
	Smtp    EmailSmtpConfig `json:"smtp"`
}

type EmailSmtpConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}
