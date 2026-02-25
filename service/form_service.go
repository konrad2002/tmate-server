package service

import (
	"github.com/konrad2002/tmate-server/model"
	"github.com/konrad2002/tmate-server/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FormService struct {
	formRepository repository.FormRepository
	historyService HistoryService
}

func NewFormService(fr repository.FormRepository, hs HistoryService) FormService {
	return FormService{
		formRepository: fr,
		historyService: hs,
	}
}

func (fs *FormService) GetAll() ([]model.Form, error) {
	return fs.formRepository.GetFormsByBsonDocument(bson.D{})
}

func (fs *FormService) GetAllForUser(userId primitive.ObjectID) ([]model.Form, error) {
	return fs.formRepository.GetFormsByBsonDocument(
		bson.D{
			{"$or",
				bson.A{
					bson.D{{"owner_user_id", userId}},
					bson.D{{"owner_user_id", primitive.NilObjectID}},
					bson.D{{"owner_user_id", bson.D{{"$exists", false}}}},
				},
			},
		})
}

func (fs *FormService) GetFormById(id primitive.ObjectID) (model.Form, error) {
	return fs.formRepository.GetFormByBsonDocument(bson.D{{"_id", id}})
}

func (fs *FormService) AddForm(form model.Form) (model.Form, error) {
	newForm, err := fs.formRepository.SaveForm(form)
	if err != nil {
		return model.Form{}, err
	}

	fs.historyService.LogFormAction(primitive.NilObjectID, model.HistoryActionCreate, newForm.Identifier)

	return newForm, nil
}

func (fs *FormService) UpdateForm(form model.Form) (model.Form, error) {
	newForm, err := fs.formRepository.UpdateForm(form)
	if err != nil {
		return model.Form{}, err
	}

	fs.historyService.LogFormAction(primitive.NilObjectID, model.HistoryActionModify, newForm.Identifier)

	return newForm, nil
}

func (fs *FormService) RemoveForm(id primitive.ObjectID) (model.Form, error) {
	return fs.formRepository.RemoveFormById(id)
}
