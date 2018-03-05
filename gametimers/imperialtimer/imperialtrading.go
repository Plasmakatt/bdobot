package imperialtimer

import (
	"time"
)

type ImperialTrading struct {
	SecondsUntilReset int64
}

func NewImperialTrading() ImperialTrading {
	return ImperialTrading{SecondsUntilReset: 4 * 60 * 60 - time.Now().Unix() % (60 * 60 * 4)}
}