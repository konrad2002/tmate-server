package service

import (
	"github.com/konrad2002/tmate-server/model"
	"github.com/konrad2002/tmate-server/repository"
	"go.mongodb.org/mongo-driver/bson"
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
