package models

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

const TimeFormat = "15:04:05.000"

type TimeString struct {
	value string
}

func NewTimeString(s string) (TimeString, error) {
	_, err := time.Parse(TimeFormat, s)
	if err != nil {
		return TimeString{}, err
	}
	return TimeString{value: s}, nil
}
func (ts *TimeString) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	tmp, err := NewTimeString(s)
	if err != nil {
		return err
	}

	*ts = tmp
	return nil
}
func NewTimeStringFromTime(t time.Time) TimeString {
	return TimeString{value: t.Format(TimeFormat)}
}

func (ts TimeString) String() string {
	return ts.value
}

func (ts TimeString) Time() (time.Time, error) {
	return time.Parse(TimeFormat, ts.value)
}

func (ts TimeString) MustTime() time.Time {
	t, _ := time.Parse(TimeFormat, ts.value)
	return t
}

func ParseDuration(durStr string) (time.Duration, error) {
	parts := strings.Split(durStr, ":")
	if len(parts) != 3 {
		return 0, errors.New("invalid duration format")
	}

	hours, _ := strconv.Atoi(parts[0])
	minutes, _ := strconv.Atoi(parts[1])
	seconds, _ := strconv.Atoi(parts[2])

	return time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second, nil
}
