package service

import (
	"github.com/konrad2002/tmate-server/misc"
	"github.com/konrad2002/tmate-server/model"
	"github.com/konrad2002/tmate-server/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type QueryService struct {
	queryRepository repository.QueryRepository
	historyService  HistoryService
}

func NewQueryService(qr repository.QueryRepository, hs HistoryService) QueryService {
	return QueryService{
		queryRepository: qr,
		historyService:  hs,
	}
}

func (qs *QueryService) GetAll() ([]model.Query, error) {
	return qs.queryRepository.GetQueriesByBsonDocument(bson.D{})
}

func (qs *QueryService) GetAllForUser(userId primitive.ObjectID) ([]model.Query, error) {
	return qs.queryRepository.GetQueriesByBsonDocument(
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

func (qs *QueryService) GetQueryById(id primitive.ObjectID) (model.Query, error) {
	return qs.queryRepository.GetQueryByBsonDocument(bson.D{{"_id", id}})
}

func (qs *QueryService) SaveExample() (model.Query, error) {
	request := bson.D{
		{"$and",
			bson.A{
				bson.D{
					{"data.trainingsgruppe",
						bson.D{
							{"$in",
								bson.A{
									"kader",
									"talente",
								},
							},
						},
					},
				},
				bson.D{{"data.dsv_lizenz_aktiv", bson.D{{"$ne", false}}}},
				bson.D{
					{"$or",
						bson.A{
							bson.D{{"data.wohnort", "Olbernhau"}},
							bson.D{{"data.wohnort", "Marienberg"}},
						},
					},
				},
			},
		},
	}

	projection := bson.D{
		{"data.vorname", 1},
		{"data.nachname", 1},
		{"data.wohnort", 1},
		{"data.aufgaben", 1},
	}

	query := model.Query{
		Name:       "Example Query",
		Filter:     request,
		Projection: projection,
		Public:     true,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	return qs.queryRepository.SaveQuery(query)
}

func (qs *QueryService) AddQuery(query model.Query) (model.Query, error) {
	query.Filter = misc.ConvertToBSOND(query.FilterJson)

	newQuery, err := qs.queryRepository.SaveQuery(query)
	if err != nil {
		return model.Query{}, err
	}

	qs.historyService.LogQueryAction(primitive.NilObjectID, model.HistoryActionCreate, newQuery.Identifier)

	return newQuery, nil
}

func (qs *QueryService) UpdateQuery(query model.Query) (model.Query, error) {
	query.Filter = misc.ConvertToBSOND(query.FilterJson)

	newQuery, err := qs.queryRepository.UpdateQuery(query)
	if err != nil {
		return model.Query{}, err
	}

	qs.historyService.LogQueryAction(primitive.NilObjectID, model.HistoryActionModify, newQuery.Identifier)

	return newQuery, nil
}

func (qs *QueryService) RemoveQuery(id primitive.ObjectID) (model.Query, error) {
	return qs.queryRepository.RemoveQueryById(id)
}
