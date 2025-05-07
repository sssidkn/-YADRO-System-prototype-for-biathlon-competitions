package controller

import (
	"biathlon-competition-system/internal/events"
	"biathlon-competition-system/internal/models"
	"biathlon-competition-system/internal/processor"
)

type Processor interface {
	ProcessEvent(comp *models.Competitor, event models.Event) models.EventResult
}

type CompetitionController struct {
	competitorsInfo map[int]*models.Competitor
	events          []models.Event
	outLogs         []string
	proc            Processor
}

type Option func(controller *CompetitionController) error

func NewCompetitionController(opts ...Option) *CompetitionController {
	ctrl := &CompetitionController{}
	for _, opt := range opts {
		err := opt(ctrl)
		if err != nil {
			return nil
		}
	}
	ctrl.outLogs = make([]string, 0, 10)
	return ctrl
}

func WithProcessor(p *processor.Processor) Option {
	return func(c *CompetitionController) error {
		c.proc = p
		return nil
	}
}

func WithEventsFile(filename string) Option {
	return func(c *CompetitionController) error {
		competitionEvents, err := events.ParseEvents(filename)
		if err != nil {
			return err
		}
		c.events = competitionEvents
		return nil
	}
}

func (c *CompetitionController) ProcessCompetition() {
	for _, event := range c.events {
		eventResult := c.proc.ProcessEvent(event.Competitor, event)
		c.outLogs = append(c.outLogs, eventResult.Msg)
		for eventResult.OutgoingEvent != nil {
			eventResult = c.proc.ProcessEvent(event.Competitor, *eventResult.OutgoingEvent)
			c.outLogs = append(c.outLogs, eventResult.Msg)
		}
	}
}

func GetOutputLog() {

}

func GetResultTable() {

}
