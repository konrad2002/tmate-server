package service

import (
	"fmt"
	"testing"
	"time"
)

type visitedDate struct {
	visited bool
	date    time.Time
}

func TestAttestService_DateRange(t *testing.T) {
	now1 := time.Date(2025, 8, 30, 14, 0, 0, 0, time.UTC)
	now2 := time.Date(2025, 8, 31, 14, 0, 0, 0, time.UTC)
	now3 := time.Date(2025, 9, 1, 14, 0, 0, 0, time.UTC)

	attests := []*visitedDate{
		{false, time.Date(2025, 9, 28, 13, 0, 0, 0, time.UTC)},
		{false, time.Date(2025, 9, 28, 14, 0, 0, 0, time.UTC)},
		{false, time.Date(2025, 9, 28, 15, 0, 0, 0, time.UTC)},
		{false, time.Date(2025, 9, 29, 13, 0, 0, 0, time.UTC)},
		{false, time.Date(2025, 9, 29, 14, 0, 0, 0, time.UTC)},
		{false, time.Date(2025, 9, 29, 15, 0, 0, 0, time.UTC)},
		{false, time.Date(2025, 9, 30, 13, 0, 0, 0, time.UTC)},
		{false, time.Date(2025, 9, 30, 14, 0, 0, 0, time.UTC)},
		{false, time.Date(2025, 9, 30, 15, 0, 0, 0, time.UTC)},
		{false, time.Date(2025, 10, 1, 13, 0, 0, 0, time.UTC)},
		{false, time.Date(2025, 10, 1, 14, 0, 0, 0, time.UTC)},
		{false, time.Date(2025, 10, 1, 15, 0, 0, 0, time.UTC)},
		{false, time.Date(2025, 10, 2, 13, 0, 0, 0, time.UTC)},
		{false, time.Date(2025, 10, 2, 14, 0, 0, 0, time.UTC)},
		{false, time.Date(2025, 10, 2, 15, 0, 0, 0, time.UTC)},
	}

	testRangeWithDate(now1, &attests)
	testRangeWithDate(now2, &attests)
	testRangeWithDate(now3, &attests)

	for _, attest := range attests {
		if attest.visited {
			fmt.Printf("[VISITED]   - %v\n", attest.date)
		} else {
			fmt.Printf("[UNVISITED] - %v\n", attest.date)
		}
	}
}

func testRangeWithDate(now time.Time, attests *[]*visitedDate) {
	println("======")
	fmt.Printf("Test with now = %v\n", now)
	println("------")

	after := now.AddDate(0, 0, +29)
	before := now.AddDate(0, 0, +30)
	fmt.Printf("after: %v | before: %v\n", after, before)

	for _, attest := range *attests {
		if attest.date.Before(before) && (attest.date.After(after) || attest.date.Equal(after)) {
			if attest.visited {
				panic("visited attest twice!!")
			}
			attest.visited = true
			fmt.Printf("attest:%v\n", attest.date)
		}
	}
}
