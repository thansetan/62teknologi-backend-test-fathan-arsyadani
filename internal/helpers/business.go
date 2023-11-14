package helpers

import (
	"time"
)

func IsOpen(currentTime time.Time, openTime, closeTime string) *bool {
	currentHour := currentTime.Format("1504")
	isOpen := (openTime <= currentHour && closeTime > currentHour)
	return &isOpen
}
