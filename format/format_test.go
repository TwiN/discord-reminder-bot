package format

import (
	"testing"
	"time"
)

func TestPrettyDuration(t *testing.T) {
	type Scenario struct {
		Name           string
		Input          time.Duration
		ExpectedOutput string
	}
	scenarios := []Scenario{
		{
			Name:           "59m",
			Input:          59 * time.Minute,
			ExpectedOutput: "59 minutes",
		},
		{
			Name:           "4m13s",
			Input:          4*time.Minute + 13*time.Second,
			ExpectedOutput: "4 minutes and 13 seconds",
		},
		{
			Name:           "1h",
			Input:          time.Hour,
			ExpectedOutput: "1 hour",
		},
		{
			Name:           "2h15m",
			Input:          2*time.Hour + 15*time.Minute,
			ExpectedOutput: "2 hours and 15 minutes",
		},
		{
			Name:           "26h30m",
			Input:          26*time.Hour + 30*time.Minute,
			ExpectedOutput: "1 day, 2 hours and 30 minutes",
		},
		{
			Name:           "26h3m4s",
			Input:          26*time.Hour + 3*time.Minute + 4*time.Second,
			ExpectedOutput: "1 day, 2 hours, 3 minutes and 4 seconds",
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			output := PrettyDuration(scenario.Input)
			if output != scenario.ExpectedOutput {
				t.Errorf("[%s] expected output to be %v, got %v", scenario.Name, scenario.ExpectedOutput, output)
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	type Scenario struct {
		Input          string
		ExpectedOutput time.Duration
		ExpectedErr    error
	}
	scenarios := []Scenario{
		{
			Input:          "1d",
			ExpectedOutput: 24 * time.Hour,
			ExpectedErr:    nil,
		},
		{
			Input:          "50d",
			ExpectedOutput: 1200 * time.Hour,
			ExpectedErr:    nil,
		},
		{
			Input:          "1h",
			ExpectedOutput: time.Hour,
			ExpectedErr:    nil,
		},
		{
			Input:          "1hr",
			ExpectedOutput: time.Hour,
			ExpectedErr:    nil,
		},
		{
			Input:          "1hour",
			ExpectedOutput: time.Hour,
			ExpectedErr:    nil,
		},
		{
			Input:          "5m",
			ExpectedOutput: 5 * time.Minute,
			ExpectedErr:    nil,
		},
		{
			Input:          "5min",
			ExpectedOutput: 5 * time.Minute,
			ExpectedErr:    nil,
		},
		{
			Input:          "5mins",
			ExpectedOutput: 5 * time.Minute,
			ExpectedErr:    nil,
		},
		{
			Input:          "5minute",
			ExpectedOutput: 5 * time.Minute,
			ExpectedErr:    nil,
		},
		{
			Input:          "5minutes",
			ExpectedOutput: 5 * time.Minute,
			ExpectedErr:    nil,
		},
		{
			Input:          "2d1h",
			ExpectedOutput: 49 * time.Hour,
			ExpectedErr:    nil,
		},
		{
			Input:          "3h2m1s",
			ExpectedOutput: 3*time.Hour + 2*time.Minute + 1*time.Second,
			ExpectedErr:    nil,
		},
		{
			Input:          "4d3h2m1s",
			ExpectedOutput: 99*time.Hour + 2*time.Minute + 1*time.Second,
			ExpectedErr:    nil,
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.Input, func(t *testing.T) {
			output, err := ParseDuration(scenario.Input)
			if err != scenario.ExpectedErr {
				t.Errorf("[%s] expected error to be %v, got %v", scenario.Input, scenario.ExpectedErr, err)
			}
			if output != scenario.ExpectedOutput {
				t.Errorf("[%s] expected output to be %v, got %v", scenario.Input, scenario.ExpectedOutput, output)
			}
		})
	}
}
