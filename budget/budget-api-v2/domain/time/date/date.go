package date

import (
	"errors"
	"time"
)

type Date struct {
	t time.Time
}

func (d *Date) GetFormattedDate() string {
	return d.t.Format("02/01/2006")
}

func (d *Date) GetIsoFormattedDate() string {
	return d.t.Format("2006-01-02")
}

func FirstDateOfMonth(month Month, year Year) (*Date, error) {
	return nil, nil
}

func LastDateOfMonth(month Month, year Year) (*Date, error) {
	return nil, nil
}

// DateFor tries several common layouts and returns a Date.
func DateFor(date string) (*Date, error) {
	layouts := []string{
		"02/01/2006",
	}
	var parseErr error
	for _, l := range layouts {
		if t, err := time.Parse(l, date); err == nil {
			return &Date{t: t}, nil
		} else {
			parseErr = err
		}
	}
	return nil, parseErr
}

// IsoDateFor parses an ISO date string like "2006-01-02" or RFC3339 datetime.
func IsoDateFor(date string) (*Date, error) {
	if t, err := time.Parse("2006-01-02", date); err == nil {
		return &Date{t: t}, nil
	}
	if t, err := time.Parse(time.RFC3339, date); err == nil {
		return &Date{t: t}, nil
	}
	return nil, errors.New("invalid ISO date format")
}
