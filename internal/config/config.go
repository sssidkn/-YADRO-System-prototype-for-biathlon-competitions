package config

import (
	"biathlon-competition-system/internal/models"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	Laps        int               `json:"laps" validate:"required"`
	LapLen      int               `json:"lapLen" validate:"required"`
	PenaltyLen  int               `json:"penaltyLen" validate:"required"`
	FiringLines int               `json:"firingLines" validate:"required"`
	Start       models.TimeString `json:"start" validate:"required"`
	StartDelta  string            `json:"startDelta" validate:"required,timeformat"`
}

type Option func(*Config) error

func New(opts ...Option) (*Config, error) {
	config := &Config{}
	for _, opt := range opts {
		if err := opt(config); err != nil {
			return nil, fmt.Errorf("config option failed: %w", err)
		}
	}
	ok, err := config.isValid()
	if !ok {
		return nil, fmt.Errorf("failed config validation: %w", err)
	}
	return config, nil
}

func (c *Config) parse(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&c); err != nil {
		return err
	}
	return nil
}

func FromFile(path string) Option {
	return func(c *Config) error {
		return c.parse(path)
	}
}

func (c *Config) isValid() (bool, error) {
	validate := validator.New()
	err := validate.RegisterValidation("timeformat", func(fl validator.FieldLevel) bool {
		_, err := time.Parse("15:04:05", fl.Field().String())
		return err == nil
	})
	if err != nil {
		return false, err
	}
	err = validate.Struct(c)
	return err == nil, err
}

func (c *Config) GetStartDeltaDuration() time.Duration {
	d, _ := models.ParseDuration(c.StartDelta)
	return d
}
