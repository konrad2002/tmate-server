package service

import (
	"github.com/konrad2002/tmate-server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AttestService struct {
	memberService MemberService
}

func NewAttestService(ms MemberService) AttestService {
	return AttestService{
		memberService: ms,
	}
}

func (as *AttestService) GetMembersWithAttestInOneMonth() ([]model.Member, error) {
	now := time.Now()

	query := model.Query{
		Identifier: primitive.NewObjectID(),
		Name:       "member-attest-in-one-month",
		Filter: bson.D{
			{"$and",
				bson.A{
					bson.D{
						{"data.attest_benoetigt", true},
						{"data.attestdatum", bson.D{{"$lte", time.Date(now.Year()-1, now.Month()+1, now.Day(), 0, 0, 0, 0, now.Location())}}},
						{"data.attestdatum", bson.D{{"$gt", time.Date(now.Year()-1, now.Month()+1, now.Day()-1, 0, 0, 0, 0, now.Location())}}},
					},
				},
			},
		},
		FilterJson:  nil,
		Projection:  nil,
		Sort:        bson.D{},
		OwnerUserId: primitive.ObjectID{},
		Public:      false,
		CreatedAt:   time.Now(),
		ModifiedAt:  time.Now(),
	}

	members, _, _, err := as.memberService.GetAllByQuery(query)
	if err != nil {
		return nil, err
	}

	return *members, nil
}

func (as *AttestService) GetMembersWithAttestOverdueToday() ([]model.Member, error) {
	now := time.Now()

	query := model.Query{
		Identifier: primitive.NewObjectID(),
		Name:       "member-attest-overdue-today",
		Filter: bson.D{
			{"$and",
				bson.A{
					bson.D{
						{"data.attest_benoetigt", true},
						{"data.attestdatum", bson.D{{"$lte", time.Date(now.Year()-1, now.Month(), now.Day(), 0, 0, 0, 0, now.Location())}}},
						{"data.attestdatum", bson.D{{"$gt", time.Date(now.Year()-1, now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())}}},
					},
				},
			},
		},
		FilterJson:  nil,
		Projection:  nil,
		Sort:        bson.D{},
		OwnerUserId: primitive.ObjectID{},
		Public:      false,
		CreatedAt:   time.Now(),
		ModifiedAt:  time.Now(),
	}

	members, _, _, err := as.memberService.GetAllByQuery(query)
	if err != nil {
		return nil, err
	}

	return *members, nil
}
