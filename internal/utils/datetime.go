package utils

import (
	"time"
)

func TodayIs() string {
	currentTime := time.Now()

	return currentTime.Format("2006-01-02")
}
