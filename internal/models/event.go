package models

type EventDefinition struct {
	ID          int
	ExtraParams string
	Comment     string
}

type Event struct {
	ID          int
	ExtraParams string
	Comment     string
	Competitors []Competitor
}
