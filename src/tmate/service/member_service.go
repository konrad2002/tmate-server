package service

import (
	"fmt"
	"github.com/konrad2002/tmate-server/repository"
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
