package processor

import (
	"biathlon-competition-system/internal/config"
	"biathlon-competition-system/internal/models"
	"fmt"
	"strconv"
	"time"
)

type EventHandler func(*models.Competitor, models.Event) models.EventResult

type Processor struct {
	handlers map[int]EventHandler
	config   config.Config
}

func New(cfg config.Config) *Processor {
	proc := &Processor{}
	proc.initHandlers()
	proc.config = cfg
	return proc
}

func (p *Processor) initHandlers() {
	p.handlers = map[int]EventHandler{
		1:  p.register,
		2:  p.setStartTime,
		3:  p.setOnStartLine,
		4:  p.start,
		5:  p.setOnFiringRange,
		6:  p.hit,
		7:  p.leftFiringRange,
		8:  p.enterPenaltyLaps,
		9:  p.leftPenaltyLaps,
		10: p.endMainLap,
		11: p.cantContinue,
		32: p.disqualify,
		33: p.finish,
	}
}

func (p *Processor) ProcessEvent(comp *models.Competitor, event models.Event) models.EventResult {
	return p.handlers[event.ID](comp, event)
}

func (p *Processor) register(comp *models.Competitor, event models.Event) models.EventResult {
	comp.Registered = true

	startTime, err := p.config.Start.Time()
	if err != nil {
		return models.EventResult{Err: err}
	}

	compStartTime := startTime.Add(p.config.GetStartDeltaDuration() * time.Duration(comp.ID-1))
	comp.StartTime = models.NewTimeStringFromTime(compStartTime)

	return models.EventResult{
		Msg: fmt.Sprintf("[%s] The competitor(%d) registered", event.TimeStamp, comp.ID),
	}
}

func (p *Processor) setStartTime(comp *models.Competitor, event models.Event) models.EventResult {
	startTime, err := time.Parse("15:04:05.000", event.ExtraParam)
	if err != nil {
		return models.EventResult{
			Err: fmt.Errorf("invalid start time format for competitor %d", comp.ID),
		}
	}
	comp.StartedAt = models.NewTimeStringFromTime(startTime)
	return models.EventResult{
		Msg: fmt.Sprintf("[%s] The start time for the competitor(%d) was set by a draw to %s",
			event.TimeStamp, comp.ID, event.ExtraParam),
	}
}

func (p *Processor) setOnStartLine(comp *models.Competitor, event models.Event) models.EventResult {
	comp.OnStartLine = true
	return models.EventResult{
		Msg: fmt.Sprintf("[%s] The competitor(%d) is on the start line", event.TimeStamp, comp.ID),
	}
}

func (p *Processor) start(comp *models.Competitor, event models.Event) models.EventResult {
	if !comp.OnStartLine {
		return models.EventResult{
			Err: fmt.Errorf("competitor %d started without being on start line", comp.ID),
		}
	}
	comp.OnStartLine = false
	comp.CurrentLap = 1

	//return models.EventResult{
	//	Msg: fmt.Sprintf("The competitor(%d) is disqualified for late start", comp.ID),
	//	OutgoingEvent: &models.Event{
	//		TimeStamp:  event.TimeStamp,
	//		ID:         32,
	//		Competitor: comp,
	//	},
	//}
	//

	return models.EventResult{
		Msg: fmt.Sprintf("The competitor(%d) has started", comp.ID),
	}
}

func (p *Processor) setOnFiringRange(comp *models.Competitor, event models.Event) models.EventResult {
	comp.OnFiringRange = true
	rangeNum, err := strconv.Atoi(event.ExtraParam)
	if err != nil {
		return models.EventResult{
			Err: fmt.Errorf("invalid firing range number: %w", err),
		}
	}
	comp.FiringRangeNum = rangeNum
	return models.EventResult{
		Msg: fmt.Sprintf("The competitor(%d) is on firing range(%d)",
			comp.ID, comp.FiringRangeNum),
	}
}

func (p *Processor) hit(comp *models.Competitor, event models.Event) models.EventResult {
	if !comp.OnFiringRange {
		return models.EventResult{
			Err: fmt.Errorf("competitor %d hit target without being on firing range", comp.ID),
		}
	}

	targetNum := event.ExtraParam
	comp.Shots++
	comp.Hits++

	return models.EventResult{
		Msg: fmt.Sprintf("The target(%d) has been hit by competitor(%d)",
			targetNum, comp.ID),
	}
}

func (p *Processor) leftFiringRange(comp *models.Competitor, event models.Event) models.EventResult {
	comp.OnFiringRange = false
	return models.EventResult{
		Msg: fmt.Sprintf("The competitor(%d) left the firing range", comp.ID),
	}
}

func (p *Processor) enterPenaltyLaps(comp *models.Competitor, event models.Event) models.EventResult {
	comp.OnPenaltyLap = true
	comp.PenaltyEnterTime = models.NewTimeStringFromTime(event.TimeStamp)
	return models.EventResult{
		Msg: fmt.Sprintf("The competitor(%d) entered penalty laps", comp.ID),
	}
}

func (p *Processor) leftPenaltyLaps(comp *models.Competitor, event models.Event) models.EventResult {
	//if comp.OnPenaltyLap {
	//	penaltyDuration := event.TimeStamp.Sub(comp.PenaltyEnterTime)
	//	comp.PenaltyTime += penaltyDuration
	//}
	comp.OnPenaltyLap = false
	return models.EventResult{
		Msg: fmt.Sprintf("The competitor(%d) left penalty laps", comp.ID),
	}
}

func (p *Processor) endMainLap(comp *models.Competitor, event models.Event) models.EventResult {
	_, err := comp.LastEventTime.Time()
	if err != nil {
		return models.EventResult{
			Err: fmt.Errorf("invalid last event time: %w", err),
		}
	}
	//lapDuration := event.TimeStamp.Sub(t)
	//comp.LapTimes = append(comp.LapTimes, models.NewTimeStringFromTime(lapDuration))
	comp.CurrentLap++
	comp.LastEventTime = models.NewTimeStringFromTime(event.TimeStamp)

	msg := fmt.Sprintf("The competitor(%d) ended main lap %d", comp.ID, comp.CurrentLap-1)

	if comp.CurrentLap > p.config.Laps {
		comp.Finished = true
		return models.EventResult{
			Msg: msg,
			OutgoingEvent: &models.Event{
				TimeStamp:  event.TimeStamp,
				ID:         33,
				Competitor: comp,
			},
		}
	}

	return models.EventResult{
		Msg: msg,
	}
}

func (p *Processor) cantContinue(comp *models.Competitor, event models.Event) models.EventResult {
	comp.CantContinue = true
	return models.EventResult{
		Msg: fmt.Sprintf("The competitor(%d) can't continue: %s",
			comp.ID, event.ExtraParam),
		OutgoingEvent: &models.Event{
			TimeStamp:  event.TimeStamp,
			ID:         34,
			Competitor: comp,
			ExtraParam: event.ExtraParam,
		},
	}
}

func (p *Processor) disqualify(comp *models.Competitor, event models.Event) models.EventResult {
	return models.EventResult{}
}

func (p *Processor) finish(comp *models.Competitor, event models.Event) models.EventResult {
	return models.EventResult{}
}
