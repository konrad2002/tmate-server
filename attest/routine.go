package attest

import (
	"fmt"
	"github.com/konrad2002/tmate-server/service"
	"time"
)

var notificationHour = 14

func StartAttestRoutine(attestService service.AttestService) {
	go func() {
		for {
			now := time.Now()

			var nextRun time.Time
			if now.Hour() >= notificationHour {
				nextRun = time.Date(now.Year(), now.Month(), now.Day()+1, notificationHour, 0, 0, 0, now.Location())
			} else {
				nextRun = time.Date(now.Year(), now.Month(), now.Day(), notificationHour, 0, 0, 0, now.Location())
			}

			duration := time.Until(nextRun)
			fmt.Println("Next run in:", duration)

			time.Sleep(duration)

			err := runTask(attestService)
			if err != nil {
				println("Attest Routine had an error:")
				println(err.Error())
			}
		}
	}()
}

func runTask(attestService service.AttestService) error {
	return attestService.RunAttestRountine()
}
