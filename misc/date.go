package misc

import (
	"fmt"
	"time"
)

var dateFormats = []string{
	"2006-01-02",
	"2006/01/02",
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05.000Z",
	"2006-01-02T15:04:05-07:00",
	"2006-01-02T15:04:05.000-07:00",
}

func ParseDate(input string) (time.Time, error) {
	for _, layout := range dateFormats {
		if t, err := time.Parse(layout, input); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported date format: %s", input)
}
