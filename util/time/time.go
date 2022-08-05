package time

import "time"

//go:generate mockery --name=Time
type Time interface {
	Now() time.Time
}

type timeAdapter struct{}

func New() Time {
	return timeAdapter{}
}

func (t timeAdapter) Now() time.Time {
	return time.Now()
}
