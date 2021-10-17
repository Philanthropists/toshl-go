package toshl

import (
	"strings"
	"time"
)

// Currency represents a Toshl supported currency
type Currency struct {
	Code  *string  `json:"code"`
	Rate  *float64 `json:"rate,omitempty"`
	Fixed *bool    `json:"fixed,omitempty"`
}

// Median represents a Toshl median
type Median struct {
	Expenses float64 `json:"expenses"`
	Incomes  float64 `json:"incomes"`
}

// Goal represents a Toshl goal
type Goal struct {
	Amount float64 `json:"amount"`
	Start  string  `json:"start"`
	End    string  `json:"end"`
}

// Recurrence represents a Toshl recurrence
type Recurrence struct {
	Frequency string `json:"frequency"`
	Interval  int    `json:"interval"`
	Start     string `json:"start"`
	Iteration int    `json:"iteration"`
}

// CategoryCounts represents a Toshl count
type CategoryCounts struct {
	Entries int `json:"entries"`
	Tags    int `json:"tags"`
}

const DateFormat = "2006-01-02"

type Date time.Time

func (v Date) MarshalJSON() ([]byte, error) {
	asTime := time.Time(v)
	return []byte("\"" + asTime.Format(DateFormat) + "\""), nil
}

func (v *Date) UnmarshalJSON(b []byte) error {
	cleaned := strings.Trim(string(b), "\"")
	timeDate, err := time.Parse(DateFormat, cleaned)
	if err != nil {
		return err
	}
	*v = Date(timeDate)
	return nil
}

func (v Date) String() string {
	asTime := time.Time(v)
	return asTime.Format(DateFormat)
}

const TimeFormat = "15:04:05"

type Time time.Time

func (v Time) MarshalJSON() ([]byte, error) {
	asTime := time.Time(v)
	return []byte(asTime.Format(TimeFormat)), nil
}

func (v *Time) UnmarshalJSON(b []byte) error {
	cleaned := strings.Trim(string(b), "\"")
	timeDate, err := time.Parse(TimeFormat, cleaned)
	if err != nil {
		return err
	}
	*v = Time(timeDate)
	return nil
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
