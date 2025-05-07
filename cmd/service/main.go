package main

import (
	"biathlon-competition-system/internal/config"
	"biathlon-competition-system/internal/controller"
	"biathlon-competition-system/internal/events"
	"biathlon-competition-system/internal/processor"
	"fmt"
)

func main() {
	cfg, err := config.New(config.FromFile("config/config.json"))
	if err != nil {
		panic(err)
	}
	eventDefinitions, err := events.ParseEventDefinitions("config/incoming events")
	if err != nil {
		panic(err)
	}
	proc := processor.New(*cfg)
	ctrl := controller.NewCompetitionController(
		controller.WithProcessor(proc),
		controller.WithEventsFile("config/events"),
	)
	ctrl.ProcessCompetition()
	fmt.Println(ctrl, cfg, eventDefinitions)
}
