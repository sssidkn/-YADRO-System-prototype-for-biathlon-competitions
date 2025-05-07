package models

type Competitor struct {
	ID               int
	StartTime        TimeString
	StartedAt        TimeString
	Disqualified     bool
	Finished         bool
	Registered       bool
	OnStartLine      bool
	CurrentLap       int
	OnFiringRange    bool
	FiringRangeNum   int
	Shots            int
	Hits             int
	OnPenaltyLap     bool
	PenaltyEnterTime TimeString
	PenaltyTime      TimeString
	LastEventTime    TimeString
	LapTimes         []TimeString
	CantContinue     bool
}
