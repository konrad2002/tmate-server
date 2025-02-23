package dto

import "github.com/konrad2002/tmate-server/model"

type QueryResultDto struct {
	Members []model.Member `json:"members"`
	Fields  []model.Field  `json:"fields"`
	Query   model.Query    `json:"query"`
}
