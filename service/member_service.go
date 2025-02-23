package service

import (
	"fmt"
	"github.com/konrad2002/tmate-server/model"
	"github.com/konrad2002/tmate-server/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MemberService struct {
	memberRepository repository.MemberRepository
	queryService     QueryService
	fieldService     FieldService
}

func NewMemberService(mr repository.MemberRepository, qs QueryService, fs FieldService) MemberService {
	return MemberService{
		memberRepository: mr,
		queryService:     qs,
		fieldService:     fs,
	}
}

func (ms *MemberService) PrintTest() string {
	fmt.Println("test")
	return "test"
}

func (ms *MemberService) GetAll() ([]model.Member, error) {
	return ms.memberRepository.GetMembersByBsonDocument(bson.D{})
}

func (ms *MemberService) GetById(id primitive.ObjectID) ([]model.Member, error) {
	return ms.memberRepository.GetMembersByBsonDocument(bson.D{{"_id", id}})
}

func (ms *MemberService) GetAllByQuery(queryId primitive.ObjectID) (*[]model.Member, *[]model.Field, *model.Query, error) {
	query, err := ms.queryService.GetQueryById(queryId)
	if err != nil {
		fmt.Println(err)
		return nil, nil, nil, err
	}

	members, err := ms.memberRepository.GetMembersByBsonDocumentWithOptions(
		query.Filter,
		options.Find().SetProjection(query.Projection).SetSort(query.Sort),
	)
	if err != nil {
		fmt.Println(err)
		return nil, nil, nil, err
	}

	fields, err := ms.fieldService.GetAllForQuery(query)
	if err != nil {
		fmt.Println(err)
		return nil, nil, nil, err
	}

	return &members, &fields, &query, nil
}
