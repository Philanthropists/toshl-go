package toshl

import (
	"errors"
	"net/url"
	"time"
)

type Entry struct {
	Id          *string   `json:"id,omitempty"`
	Amount      float64   `json:"amount"`
	Currency    Currency  `json:"currency"`
	Date        string    `json:"date"`
	Description *string   `json:"desc,omitempty"`
	Account     string    `json:"account"`
	Category    string    `json:"category"`
	Tags        []string  `json:"tags,omitempty"`
	Location    *Location `json:"location,omitempty"`
	Created     time.Time `json:"created"`
	Modified    string    `json:"modified"`
	Repeat      *Repeat   `json:"repeat,omitempty"`
}

type EntryQueryParams struct {
	From Date
	To   Date
}

func (a *EntryQueryParams) getQueryString() (string, error) {
	v := url.Values{}

	var nilDate Date
	var errMsg string

	if a.From == nilDate {
		errMsg = errMsg + "'from' field is mandatory;"
	}

	v.Set("from", a.From.String())

	if a.To == nilDate {
		errMsg = errMsg + "'to' field is mandatory;"
	}

	v.Set("to", a.To.String())

	if errMsg != "" {
		return "", errors.New(errMsg)
	}

	return v.Encode(), nil
}

type RepeatFrequency string

const (
	Daily   RepeatFrequency = "daily"
	Weekly  RepeatFrequency = "weekly"
	Monthly RepeatFrequency = "monthly"
	Yearly  RepeatFrequency = "yearly"
)

type RepeatType string

const (
	Automatic RepeatType = "automatic"
	Confirm   RepeatType = "confirm"
	Confirmed RepeatType = "confirmed"
)

type Repeat struct {
	Start      Date            `json:"start"`
	End        Date            `json:"end"`
	Frequency  RepeatFrequency `json:"frequency"`
	Interval   uint            `json:"interval"`
	Count      uint            `json:"count"`
	ByDay      string          `json:"byday"`
	ByMonthDay string          `json:"bymonthday"`
	BySetPos   string          `json:"bysetpos"`
	Iteration  uint            `json:"iteration"`
	IsTemplate bool            `json:"template"`
	Entries    []string        `json:"entries"`
	Type       RepeatType      `json:"type"`
}

type Transaction struct {
	Id        string                 `json:"id"`
	Amount    float64                `json:"amount"`
	Account   string                 `json:"account,omitempty"`
	Currency  Currency               `json:"currency,omitempty"`
	Images    []Image                `json:"images"`
	Reminders []Reminder             `json:"reminders"`
	Import    Import                 `json:"import"`
	Review    Review                 `json:"review"`
	Settle    Settle                 `json:"settle"`
	Split     Split                  `json:"split"`
	Readonly  []string               `json:"readonly"`
	Completed bool                   `json:"completed"`
	Deleted   bool                   `json:"deleted"`
	Extra     map[string]interface{} `json:"extra"`
}

type Image struct {
	Id       string `json:"id"`
	Path     string `json:"path"`
	Filename string `json:"filename"`
	Type     string `json:"type"`
	Status   string `json:"status"`
}

type Reminder struct {
	Period string `json:"period"`
	Number uint   `json:"number"`
	//At     Time   `json:"at"`
}

type Import struct {
	// TODO fill
}

type Review struct {
	Id        string `json:"id"`
	Type      string `json:"type"`
	Completed bool   `json:"completed"`
}

type Settle struct {
	Id string `json:"id"`
}

type Split struct {
	Parent   string   `json:"parent"`
	Children []string `json:"children"`
}
