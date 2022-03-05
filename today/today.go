package today

import (
	"time"
)

const DayFormat = "2006-01-02"

type Today struct {
	time.Time
}

func New() *Today {
	today := &Today{time.Now()}
	// Remove daytime for comparison
	// Skipping error is ok because used variable is type safe
	today, _ = Parse(today.GetString())

	return today
}

func (t *Today) GetDayPlus(days int) string {
	return t.AddDate(0, 0, days).Format(DayFormat)
}

func (t *Today) GetString() string {
	return t.Format(DayFormat)
}

func Parse(date string) (*Today, error) {
	parsed, err := time.Parse(DayFormat, date)

	return &Today{parsed}, err
}
