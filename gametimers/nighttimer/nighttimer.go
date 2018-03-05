package nighttimer

import (
	"time"
)

const DaytimeLength = 12000
const NighttimeLength = 2400
const DayLength = DaytimeLength + NighttimeLength

//May be used in the future, all commented out is coupled
/* type NightTimer struct {
	SecsUntilNightEnd, SecsUntilNightStart, GameHour, GameMinute int64
	IsDay bool
} */

type NightTimer struct {
	SecondsUntilNightEnd, SecondsUntilNightStart int64
	IsDay bool
}

//May be used in the future, all commented out is coupled
/* func GetGameHour(pctOfNightDone int64) int64 {
	gameHour := 9 * pctOfNightDone
	if gameHour < 2 {
		gameHour = 22 + gameHour
	} else {
		gameHour = gameHour - 2
	}
	return gameHour
} */

func GetTimerDuringNight(secsIntoGameDay int64) NightTimer {
	secsIntoGameNight := secsIntoGameDay - DaytimeLength
	secsUntilNightEnd := NighttimeLength - secsIntoGameNight
	secsUntilNightStart := secsUntilNightEnd + DaytimeLength

	//May be used in the future, all commented out is coupled
	/* 	pctOfNightDone := secsIntoGameNight / NighttimeLength
	gameHour := GetGameHour(pctOfNightDone)
	inGameHour := gameHour / 1 >> 0
	inGameMinute := gameHour % 1 * 60 >> 0 */
	return NightTimer{SecondsUntilNightEnd: secsUntilNightEnd, SecondsUntilNightStart: secsUntilNightStart, IsDay: false}
}

func GetTimerDuringDay(secsIntoGameDay int64) NightTimer {
	secsUntilNightStart := DaytimeLength - secsIntoGameDay
	secsUntilNightEnd := secsUntilNightStart + NighttimeLength
	//May be used in the future, all commented out is coupled
	/* 	pctOfDayDone := secsIntoGameDay / DaytimeLength
	gameHour := 7 + (22 - 7) * pctOfDayDone // Night starts at 22 and ends at 7
	inGameHour := gameHour / 1 >> 0
	inGameMinute := gameHour % 1 * 60 >> 0 */
	return NightTimer{SecondsUntilNightEnd: secsUntilNightEnd, SecondsUntilNightStart: secsUntilNightStart, IsDay: true}
}

func New() NightTimer {
	var timer NightTimer
	secsIntoGameDay := (time.Now().Unix() + DaytimeLength + 20 * 60) % (DayLength)
	if secsIntoGameDay >= DaytimeLength {
		timer = GetTimerDuringNight(secsIntoGameDay)
	} else {
		timer = GetTimerDuringDay(secsIntoGameDay)
	}
	return timer
}