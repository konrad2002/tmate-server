package service

import (
	"fmt"
	"github.com/konrad2002/tmate-server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AttestService struct {
	memberService MemberService
	fieldService  FieldService
	configService ConfigService
	emailService  EmailService
}

func NewAttestService(ms MemberService, fs FieldService, cs ConfigService, es EmailService) AttestService {
	return AttestService{
		memberService: ms,
		fieldService:  fs,
		configService: cs,
		emailService:  es,
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
						{"data.attestdatum", bson.D{{"$lt", now.AddDate(-1, 0, +30)}}},  // 31.9. => 1.10.
						{"data.attestdatum", bson.D{{"$gte", now.AddDate(-1, 0, +29)}}}, // 30.9.
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
						{"data.attestdatum", bson.D{{"$lt", now.AddDate(-1, 0, 1)}}},  // 2.9.
						{"data.attestdatum", bson.D{{"$gte", now.AddDate(-1, 0, 0)}}}, // 1.9.
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

func (as *AttestService) RunAttestRountine() error {
	fmt.Println("\033[1;36m  --===[ ATTEST ROUTINE ]===-- \033[0m Task executed at:", time.Now())

	specialFields, err := as.configService.GetSpecialFields()
	if err != nil {
		return err
	}

	// find members with attest decay in one month:
	members, err := as.GetMembersWithAttestInOneMonth()
	if err != nil {
		return err
	}

	println("\033[36mMembers with attest being outdated in <1 month:\033[0m")
	for _, member := range members {
		firstName := member.Data[specialFields.FirstName].(string)
		lastName := member.Data[specialFields.LastName].(string)
		email := member.Data[specialFields.EMail].(string)
		date := (member.Data[specialFields.AttestDate]).(primitive.DateTime).Time().AddDate(1, 0, 0).Format("02.01.2006")
		fmt.Printf("%s, %s, %s, %s\n", firstName, lastName, email, date)

		// notify about attest in one month (send email)
		err := as.emailService.SendAttestEmail(firstName, lastName, email, date, "Ärztliches Attest bald ungültig", "assets/templates/attest_email_warning.html", member)
		if err != nil {
			fmt.Printf("\033[37mfailed to send mail: %s\033[0m\n", err)
		}
	}

	// find member with attest decay today:
	members2, err := as.GetMembersWithAttestOverdueToday()
	if err != nil {
		return err
	}

	println("\033[36mMembers with attest being outdated today:\033[0m")
	for _, member := range members2 {
		firstName := member.Data[specialFields.FirstName].(string)
		lastName := member.Data[specialFields.LastName].(string)
		email := member.Data[specialFields.EMail].(string)
		date := (member.Data[specialFields.AttestDate]).(primitive.DateTime).Time().AddDate(1, 0, 0).Format("02.01.2006")
		fmt.Printf("%s, %s, %s, %s\n", firstName, lastName, email, date)

		// notify about attest missing (send email)
		err := as.emailService.SendAttestEmail(firstName, lastName, email, date, "Ärztliches Attest ungültig!", "assets/templates/attest_email_missing.html", member)
		if err != nil {
			fmt.Printf("\033[37mfailed to send mail: %s\033[0m\n", err)
		}
	}

	return nil
}
