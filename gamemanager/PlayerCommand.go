package gamemanager

import (
	"time"
)

type PlayerCommand struct {
	message   string
	timestamp time.Time
}
