package service

import (
	"fmt"
	"github.com/konrad2002/tmate-server/model"
	"github.com/konrad2002/tmate-server/repository"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

type FieldService struct {
	fieldRepository repository.FieldRepository
}

func NewFieldService(mr repository.FieldRepository) FieldService {
	return FieldService{
		fieldRepository: mr,
	}
}

func (fs *FieldService) GetAll() ([]model.Field, error) {
	return fs.fieldRepository.GetFieldsByBsonDocument(bson.D{})
}

func (fs *FieldService) GetFirstFieldWithType(fieldType model.FieldType) (model.Field, error) {
	return fs.fieldRepository.GetFieldByBsonDocument(bson.D{{"type", fieldType}})
}

func (fs *FieldService) GetAllForQuery(query model.Query) ([]model.Field, error) {
	var fieldNames []string
	for _, e := range query.Projection {
		if strings.Contains(e.Key, "data.") {
			fieldNames = append(fieldNames, strings.Replace(e.Key, "data.", "", 1))
		}
	}

	fmt.Printf("field names: %s\n", fieldNames)

	return fs.fieldRepository.GetFieldsByBsonDocument(
		bson.D{
			{
				"name",
				bson.M{"$in": fieldNames},
			},
		},
	)
}

func (fs *FieldService) AddField(field model.Field) (model.Field, error) {
	return fs.fieldRepository.SaveField(field)
}
