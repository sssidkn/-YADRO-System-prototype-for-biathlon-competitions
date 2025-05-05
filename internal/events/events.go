package events

import (
	"biathlon-competition-system/internal/models"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var EventDefinitionsNumber = 11

func ParseEventDefinitions(filename string) (map[int]models.EventDefinition, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	definitions := make(map[int]models.EventDefinition, EventDefinitionsNumber)
	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		_ = scanner.Text()
	}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 3 {
			continue
		}

		idStr := strings.TrimSpace(parts[0])
		extraParams := strings.TrimSpace(parts[1])
		comment := strings.TrimSpace(parts[2])

		var id int
		id, err = strconv.Atoi(idStr)
		if err != nil {
			return nil, err
		}

		definitions[id] = models.EventDefinition{
			ID:          id,
			ExtraParams: extraParams,
			Comment:     comment,
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return definitions, nil
}

func ParseEvents(filename string) ([]models.Event, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var events []models.Event
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		event, err := parseEventLine(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err)
		}

		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return events, nil
}

// parseEventLine парсит одну строку с событием
func parseEventLine(line string) (models.Event, error) {
	parts := strings.Split(line, " ")
	if len(parts) < 3 {
		return models.Event{}, fmt.Errorf("invalid event format, expected at least 3 parts")
	}

	timeStr := parts[0]
	t, err := time.Parse("[15:04:05.000]", timeStr)
	if err != nil {
		return models.Event{}, fmt.Errorf("invalid time format: %w", err)
	}

	eventID, err := strconv.Atoi(parts[1])
	if err != nil {
		return models.Event{}, fmt.Errorf("invalid event ID: %w", err)
	}

	competitorID, err := strconv.Atoi(parts[2])
	if err != nil {
		return models.Event{}, fmt.Errorf("invalid competitor ID: %w", err)
	}

	var extraParams string
	if len(parts) > 3 {
		extraParams = strings.Join(parts[3:], " ")
	}

	return models.Event{
		ID:          eventID,
		ExtraParams: extraParams,
		Competitors: []models.Competitor{{ID: uint(competitorID), StartTime: t.Format("15:04:05")}},
	}, nil
}
