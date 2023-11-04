package utils

import (
	"log"
	"time"
)

func IsDateValue(stringDate string) (time.Time, bool) {
	timer, err := time.Parse("2006/01/02", stringDate)
	log.Print(err)
	return timer, err == nil
}
