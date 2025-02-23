package service

import (
	"fmt"
	"github.com/konrad2002/tmate-server/model"
	"github.com/konrad2002/tmate-server/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemberService struct {
	memberRepository repository.MemberRepository
	queryService     QueryService
}

func NewMemberService(mr repository.MemberRepository, qs QueryService) MemberService {
	return MemberService{
		memberRepository: mr,
		queryService:     qs,
	}
}

func (ms *MemberService) PrintTest() string {
	fmt.Println("test")
	return "test"
}

func (ms *MemberService) GetAll() ([]model.Member, error) {
	return ms.memberRepository.GetMembersByBsonDocument(bson.D{})
}

func (ms *MemberService) GetAllByQuery(queryId primitive.ObjectID) ([]model.Member, error) {
	query, err := ms.queryService.GetQueryById(queryId)
	if err != nil {
		fmt.Println(err)
		return []model.Member{}, err
	}
	return ms.memberRepository.GetMembersByBsonDocument(query.Query)
}
