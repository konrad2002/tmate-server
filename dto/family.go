package dto

import "github.com/konrad2002/tmate-server/model"

type FamilyListDto struct {
	Families map[int][]model.Member `json:"families"`
}
