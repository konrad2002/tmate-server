package service

import (
	"fmt"
	"github.com/konrad2002/tmate-server/model"
	"github.com/konrad2002/tmate-server/repository"
	"go.mongodb.org/mongo-driver/bson"
)

type MemberService struct {
	memberRepository repository.MemberRepository
}

func NewMemberService(mr repository.MemberRepository) MemberService {
	return MemberService{
		memberRepository: mr,
	}
}

func (ms *MemberService) PrintTest() string {
	fmt.Println("test")
	return "test"
}

func (ms *MemberService) GetAll() ([]model.Member, error) {
	return ms.memberRepository.GetMembersByBsonDocument(bson.D{})
}
