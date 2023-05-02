package main

import (
	"time"
)

type clock interface {
	Time() time.Time
}

type defaultClock int

func (defaultClock) Time() time.Time { return time.Now() }

var defaultClockInstance defaultClock
