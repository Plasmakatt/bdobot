package timeconversion

import (
	"strings"
	"strconv"
	"time"
)

func GetHMSFromSeconds(seconds int64) string {
	hours, remainingSecs := ConvertSecondsToHours(seconds)
	minutes, remainingSecs := ConvertSecondsToMinutes(remainingSecs)
	var messages []string
	messages = append(messages, GetHoursAsString(hours))
	minutesAsString := GetMinutesAsString(minutes)
	if len(minutesAsString) > 0  {
		messages = append(messages, minutesAsString + " and")
	}
	messages = append(messages, GetSecondsAsString(remainingSecs))
	return strings.Join(messages, " ")
}

func ConvertSecondsToMinutes(seconds int64) (int64, int64) {
	return seconds / 60, seconds % 60
}

func ConvertSecondsToHours(seconds int64) (int64, int64) {
	return seconds / 3600, seconds % 3600
}

func GetHoursAsString(hours int64) string {
	s := ""
	if hours > 1 {
		s = strconv.FormatInt(hours, 10) + " hours"
	} else if hours == 1 {
		s =  strconv.FormatInt(hours, 10) + " hour"
	}
	return s
}

func GetMinutesAsString(minutes int64) string {
	s := ""
	if minutes > 1 {
		s = strconv.FormatInt(minutes, 10) + " minutes"
	} else if minutes == 1 {
		s = strconv.FormatInt(minutes, 10) + " minute"
	}
	return s
}

func GetSecondsAsString(seconds int64) string {
	s := ""
	if seconds > 1 {
		s = strconv.FormatInt(seconds, 10) + " seconds"
	} else if seconds == 1 {
		s = "and " + strconv.FormatInt(seconds, 10) + " second"
	}
	return s
}

func GetDateAfterSeconds(seconds int64) string {
	return time.Now().Local().Add(time.Second * time.Duration(seconds)).String()
}