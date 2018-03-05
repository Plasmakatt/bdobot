package energytimer

import (
	"strconv"
)

const EnergyPerMinute = 0.33

type EnergyTimer struct {
	CurrentEnergy, MaxEnergy string
}

func (enTimer EnergyTimer) GetRemainingSeconds() int64 {
	maxEnInInt, _ := strconv.Atoi(enTimer.MaxEnergy)
	currEnInInt, _ := strconv.Atoi(enTimer.CurrentEnergy)
	energyToFill := maxEnInInt - currEnInInt
	minutesRemaining := (float64(energyToFill) / EnergyPerMinute)
	return int64(minutesRemaining) * 60
}