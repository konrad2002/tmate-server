package attest

import (
	"fmt"
	"github.com/konrad2002/tmate-server/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func StartAttestRoutine(attestService service.AttestService, fieldService service.FieldService, configService service.ConfigService, emailService service.EmailService) {
	go func() {
		for {
			err := runTask(attestService, fieldService, configService, emailService)
			if err != nil {
				println("Attest Routine had an error:")
				println(err.Error())
			}

			now := time.Now()
			nextRun := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()+10, 0, 0, now.Location())
			duration := time.Until(nextRun)
			fmt.Println("Next run in:", duration)

			time.Sleep(duration)

		}
	}()
}

func runTask(attestService service.AttestService, fieldService service.FieldService, configService service.ConfigService, emailService service.EmailService) error {
	fmt.Println("\033[1;36m  --===[ ATTEST ROUTINE ]===-- \033[0m Task executed at:", time.Now())

	specialFields, err := configService.GetSpecialFields()
	if err != nil {
		return err
	}

	// find members with attest decay in one month:
	members, err := attestService.GetMembersWithAttestInOneMonth()
	if err != nil {
		return err
	}

	println("\033[36mMembers with attest being outdated in <1 month:\033[0m")
	for _, member := range members {
		firstName := member.Data[specialFields.FirstName].(string)
		lastName := member.Data[specialFields.LastName].(string)
		email := member.Data[specialFields.EMail].(string)
		date := (member.Data[specialFields.AttestDate]).(primitive.DateTime).Time().Format("02.01.2006")
		fmt.Printf("%s, %s, %s, %s\n", firstName, lastName, email, date)

		// notify about attest in one month (send email)
		err := emailService.SendAttestEmail(firstName, lastName, email, date, "Ärztliches Attest bald ungültig", "assets/templates/attest_email_warning.html")
		if err != nil {
			fmt.Printf("\033[37mfailed to send mail: %s\033[0m\n", err)
		}
	}

	// find member with attest decay today:
	members2, err := attestService.GetMembersWithAttestOverdueToday()
	if err != nil {
		return err
	}

	println("\033[36mMembers with attest being outdated today:\033[0m")
	for _, member := range members2 {
		firstName := member.Data[specialFields.FirstName].(string)
		lastName := member.Data[specialFields.LastName].(string)
		email := member.Data[specialFields.EMail].(string)
		date := (member.Data[specialFields.AttestDate]).(primitive.DateTime).Time().Format("02.01.2006")
		fmt.Printf("%s, %s, %s, %s\n", firstName, lastName, email, date)

		// notify about attest missing (send email)
		err := emailService.SendAttestEmail(firstName, lastName, email, date, "Ärztliches Attest ungültig!", "assets/templates/attest_email_missing.html")
		if err != nil {
			fmt.Printf("\033[37mfailed to send mail: %s\033[0m\n", err)
		}
	}

	return nil
}
