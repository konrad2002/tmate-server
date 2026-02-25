package service

import (
	"github.com/konrad2002/tmate-server/model"
	"github.com/konrad2002/tmate-server/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HistoryService struct {
	historyRepository repository.HistoryRepository
}

func NewHistoryService(hr repository.HistoryRepository) HistoryService {
	return HistoryService{
		historyRepository: hr,
	}
}

func (hs *HistoryService) GetAll() ([]model.History, error) {
	return hs.historyRepository.GetHistoriesByBsonDocument(bson.D{})
}

func (hs *HistoryService) LogMemberAction(user primitive.ObjectID, action model.HistoryAction, member primitive.ObjectID) {
	history := model.History{
		UserId:         user,
		Action:         action,
		TargetMemberId: member,
	}
	_, err := hs.historyRepository.SaveHistory(history)
	if err != nil {
		println("HISTORY WRITE ERROR: " + err.Error())
	}
}

func (hs *HistoryService) LogQueryAction(user primitive.ObjectID, action model.HistoryAction, query primitive.ObjectID) {
	history := model.History{
		UserId:        user,
		Action:        action,
		TargetQueryId: query,
	}
	_, err := hs.historyRepository.SaveHistory(history)
	if err != nil {
		println("HISTORY WRITE ERROR: " + err.Error())
	}
}

func (hs *HistoryService) LogFormAction(user primitive.ObjectID, action model.HistoryAction, form primitive.ObjectID) {
	history := model.History{
		UserId:       user,
		Action:       action,
		TargetFormId: form,
	}
	_, err := hs.historyRepository.SaveHistory(history)
	if err != nil {
		println("HISTORY WRITE ERROR: " + err.Error())
	}
}

func (hs *HistoryService) LogEMailAction(user primitive.ObjectID, member primitive.ObjectID, content string) {
	history := model.History{
		UserId:         user,
		Action:         model.HistoryActionSend,
		TargetMemberId: member,
		EMailContent:   content,
	}
	_, err := hs.historyRepository.SaveHistory(history)
	if err != nil {
		println("HISTORY WRITE ERROR: " + err.Error())
	}
}

func (hs *HistoryService) SaveHistory(history model.History) (model.History, error) {
	return hs.historyRepository.SaveHistory(history)
}
