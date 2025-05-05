package main

import (
	"biathlon-competition-system/internal/config"
	"biathlon-competition-system/internal/events"
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
	competitionEvents, err := events.ParseEvents("config/events")
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg, eventDefinitions, competitionEvents)
}
