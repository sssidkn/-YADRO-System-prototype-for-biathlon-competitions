package models

import "time"

type EventDefinition struct {
	ID          int
	ExtraParams string
	Comment     string
}

type Event struct {
	ID         int
	ExtraParam string
	Comment    string
	TimeStamp  time.Time
	Competitor *Competitor
}

type EventResult struct {
	Err           error
	Msg           string
	OutgoingEvent *Event
}
