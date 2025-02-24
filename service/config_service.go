package service

import (
	"encoding/json"
	"github.com/konrad2002/tmate-server/model"
	"os"
)

var configDir = "config/"
var specialFieldsFile = configDir + "special_fields.json"
var configFile = configDir + "config.json"

type ConfigService struct {
}

func NewConfigService() ConfigService {
	return ConfigService{}
}

func (cs *ConfigService) InitConfig() error {
	specialFields := model.SpecialFields{
		FirstName: "firstname",
		LastName:  "lastname",
		EMail:     "email",
		EMail2:    "email2",
		Address: model.AddressFields{
			Street:     "street",
			Number:     "number",
			City:       "city",
			PostalCode: "postal_code",
		},
	}

	specialFieldsJson, err := json.Marshal(specialFields)
	if err != nil {
		return err
	}

	err = os.WriteFile(specialFieldsFile, specialFieldsJson, 0644)
	if err != nil {
		return err
	}

	config := model.Config{}

	configJson, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFile, configJson, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (cs *ConfigService) GetSpecialFields() (*model.SpecialFields, error) {
	jsonString, err := os.ReadFile(specialFieldsFile)
	if err != nil {
		return nil, err
	}

	var specialFields model.SpecialFields

	err = json.Unmarshal(jsonString, &specialFields)
	if err != nil {
		return nil, err
	}

	return &specialFields, nil
}

func (cs *ConfigService) GetConfig() (*model.Config, error) {
	jsonString, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config model.Config

	err = json.Unmarshal(jsonString, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
