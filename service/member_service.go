package service

import (
	"errors"
	"fmt"
	"github.com/konrad2002/tmate-server/dto"
	"github.com/konrad2002/tmate-server/misc"
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
	configService    ConfigService
}

func NewMemberService(mr repository.MemberRepository, qs QueryService, fs FieldService, cs ConfigService) MemberService {
	return MemberService{
		memberRepository: mr,
		queryService:     qs,
		fieldService:     fs,
		configService:    cs,
	}
}

func (ms *MemberService) GetSlimMemberOptions() (*options.FindOptions, error) {

	specialFields, err := ms.configService.GetSpecialFields()
	if err != nil {
		return nil, err
	}

	queryOptions := options.FindOptions{}
	queryOptions.SetProjection(bson.D{
		{"data." + specialFields.FirstName, 1},
		{"data." + specialFields.LastName, 1},
		{"data." + specialFields.EMail, 1},
		{"data." + specialFields.Family, 1},
	})

	return &queryOptions, nil
}

func (ms *MemberService) GetAll() ([]model.Member, error) {
	return ms.memberRepository.GetMembersByBsonDocument(bson.D{})
}

func (ms *MemberService) GetById(id primitive.ObjectID) (model.Member, error) {
	return ms.memberRepository.GetMemberByBsonDocument(bson.D{{"_id", id}})
}

func (ms *MemberService) GetFamilies() (*dto.FamilyListDto, error) {
	slimMemberOptions, err := ms.GetSlimMemberOptions()
	if err != nil {
		return nil, err
	}

	members, err := ms.memberRepository.GetMembersByBsonDocumentWithOptions(bson.D{}, slimMemberOptions)
	if err != nil {
		return nil, err
	}

	familyField, err := ms.fieldService.GetFirstFieldWithType(model.Family)
	if err != nil {
		return nil, err
	}

	if familyField.Name == "" {
		err := errors.New(fmt.Sprintf("no family field defined!"))
		return nil, err
	}

	var families dto.FamilyListDto
	families.Families = make(map[int][]model.Member)

	for _, member := range members {
		familyIdValue := member.Data[familyField.Name]
		if familyIdValue != nil {
			familyId, err := misc.AnyToInt(familyIdValue)
			if err != nil {
				err := errors.New(fmt.Sprintf("failed to parse family id (%d) for member %s", member.Data[familyField.Name], member.Identifier))
				return nil, err
			}

			families.Families[familyId] = append(families.Families[familyId], member)
		}
	}

	return &families, nil
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

func (ms *MemberService) AddMember(member model.Member) (model.Member, error) {
	return ms.memberRepository.SaveMember(member)
}

func (ms *MemberService) UpdateMember(member model.Member) (model.Member, error) {
	return ms.memberRepository.UpdateMember(member)
}
