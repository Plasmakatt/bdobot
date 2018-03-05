package imperialtimer

import (
	"time"
)

type ImperialCooking struct {
	SecondsUntilReset int64
}

func NewImperialCooking() ImperialCooking {
	return ImperialCooking{SecondsUntilReset: 3 * 60 * 60 - time.Now().Unix() % (60 * 60 * 3)}
}