package service

import (
	"github.com/konrad2002/tmate-server/model"
	"github.com/konrad2002/tmate-server/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type QueryService struct {
	queryRepository repository.QueryRepository
}

func NewQueryService(qr repository.QueryRepository) QueryService {
	return QueryService{
		queryRepository: qr,
	}
}

func (qs *QueryService) GetAll() ([]model.Query, error) {
	return qs.queryRepository.GetQueriesByBsonDocument(bson.D{})
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
	return qs.queryRepository.SaveQuery(query)
}

func (qs *QueryService) UpdateQuery(query model.Query) (model.Query, error) {
	return qs.queryRepository.UpdateQuery(query)
}

func (qs *QueryService) RemoveQuery(id primitive.ObjectID) (model.Query, error) {
	return qs.queryRepository.RemoveQueryById(id)
}
