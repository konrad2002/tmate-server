package dto

import "github.com/konrad2002/tmate-server/model"

type FamilyListDto struct {
	Families map[int]*Family `json:"families"`
}

type Family struct {
	MemberCount int            `json:"member_count"`
	LastName    string         `json:"last_name"`
	Members     []model.Member `json:"members"`
}
