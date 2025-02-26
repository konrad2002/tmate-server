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
	"time"
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

	specialFields, err := ms.configService.GetSpecialFields()
	if err != nil {
		return nil, err
	}

	var families dto.FamilyListDto
	families.Families = make(map[int]*dto.Family)

	// iterate over all members, add them to their family
	for _, member := range members {
		familyIdValue := member.Data[familyField.Name]
		if familyIdValue != nil {
			familyId, err := misc.AnyToInt(familyIdValue)
			if err != nil {
				err := errors.New(fmt.Sprintf("failed to parse family id (%d) for member %s", member.Data[familyField.Name], member.Identifier))
				return nil, err
			}

			if familyId == 0 {
				continue
			}

			if families.Families[familyId] == nil {
				families.Families[familyId] = &dto.Family{
					MemberCount: 0,
					LastName:    member.Data[specialFields.LastName].(string),
				}
			}

			families.Families[familyId].MemberCount++
			families.Families[familyId].Members = append(families.Families[familyId].Members, member)
		}
	}

	// collect last names and find the highest occurrence
	for _, family := range families.Families {
		lastNames := make(map[string]int)
		for _, member := range family.Members {
			lastNames[member.Data[specialFields.LastName].(string)]++
		}

		highestCount := 0
		for s, i := range lastNames {
			if i > highestCount {
				family.LastName = s
				highestCount = i
			}
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
		fmt.Println("failed to retrieve members")
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

func (ms *MemberService) AddMember(member model.Member, familyMemberId primitive.ObjectID) (model.Member, error) {
	if !familyMemberId.IsZero() {
		var err error
		member, err = ms.createOrAddFamily(member, familyMemberId)
		if err != nil {
			return model.Member{}, err
		}
	}

	return ms.memberRepository.SaveMember(member)
}

// getFamilyIdFromMember returns the family members id;
// boolean, if member has a family id
// error: no family field; failed to parse family id
func (ms *MemberService) getFamilyIdFromMember(member model.Member, familyFieldName string) (int, bool, error) {
	if familyFieldName == "" {

		familyField, err := ms.fieldService.GetFirstFieldWithType(model.Family)
		if err != nil {
			return 0, false, err
		}

		if familyField.Name == "" {
			err := errors.New(fmt.Sprintf("no family field defined!"))
			return 0, false, err
		}

		familyFieldName = familyField.Name
	}

	familyIdValue := member.Data[familyFieldName]
	if familyIdValue != nil {
		familyId, err := misc.AnyToInt(familyIdValue)
		if err != nil {
			err := errors.New(fmt.Sprintf("failed to parse family id (%d) for member %s", member.Data[familyFieldName], member.Identifier))
			return 0, false, err
		}
		return familyId, true, nil
	} else {
		return 0, false, nil
	}

}

func (ms *MemberService) UpdateMember(member model.Member, familyMemberId primitive.ObjectID) (model.Member, error) {
	if !familyMemberId.IsZero() {
		var err error
		member, err = ms.createOrAddFamily(member, familyMemberId)
		if err != nil {
			return model.Member{}, err
		}
	}

	return ms.memberRepository.UpdateMember(member)
}

func (ms *MemberService) createOrAddFamily(member model.Member, familyMemberId primitive.ObjectID) (model.Member, error) {
	// get family field
	familyField, err := ms.fieldService.GetFirstFieldWithType(model.Family)
	if err != nil {
		return model.Member{}, err
	}

	if familyField.Name == "" {
		err := errors.New(fmt.Sprintf("no family field defined!"))
		return model.Member{}, err
	}

	// get member family
	familyMember, err := ms.GetById(familyMemberId)
	if err != nil {
		println(err.Error())
		err2 := errors.New(fmt.Sprintf("failed to lookup family member: %s\n", familyMemberId))
		return model.Member{}, err2
	}

	familyId, has, err := ms.getFamilyIdFromMember(familyMember, familyField.Name)
	if err != nil {
		return model.Member{}, err
	}

	// if no family, create one
	if !has {
		// generate id
		familyId = int(time.Now().UnixNano())

		// assign to family member
		familyMember.Data[familyField.Name] = familyId

		// save family member
		_, err := ms.memberRepository.UpdateMember(familyMember)
		if err != nil {
			err2 := errors.New(fmt.Sprintf("failed to save family member (%s): %e\n", familyMemberId, err))
			return model.Member{}, err2
		}
	}

	// save number
	member.Data[familyField.Name] = familyId

	return member, nil
}
