package date

import (
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

func TestDateIsFormattedWithDefaultFormatter(t *testing.T) {
	expectedFormattedDate := "25/02/2018"
	anotherExpectedFormattedDate := "25/03/2018"
	anotherExpectedFormattedDate2 := "25/05/2018"

	date, _ := IsoDateFor("2018-02-25")
	anotherDate, _ := IsoDateFor("2018-03-25")
	anotherDate2, _ := IsoDateFor("2018-05-25")

	assert.Equal(t, date.GetFormattedDate(), expectedFormattedDate)
	assert.Equal(t, anotherDate.GetFormattedDate(), anotherExpectedFormattedDate)
	assert.Equal(t, anotherDate2.GetFormattedDate(), anotherExpectedFormattedDate2)
}

func TestDateIsFormattedWithIsoFormatter(t *testing.T) {
	expectedFormattedDate := "2018-02-25"
	anotherExpectedFormattedDate := "2018-03-25"
	anotherExpectedFormattedDate2 := "2018-05-25"

	date, _ := IsoDateFor("2018-02-25")
	anotherDate, _ := IsoDateFor("2018-03-25")
	anotherDate2, _ := IsoDateFor("2018-05-25")

	assert.Equal(t, date.GetIsoFormattedDate(), expectedFormattedDate)
	assert.Equal(t, anotherDate.GetIsoFormattedDate(), anotherExpectedFormattedDate)
	assert.Equal(t, anotherDate2.GetIsoFormattedDate(), anotherExpectedFormattedDate2)
}

func TestFirstDateOfTheMouth(t *testing.T) {
	expectedFormattedDate := "01/04/2018"

	month := APRIL()
	year := NewYear(2018)

	firstDate, _ := FirstDateOfMonth(month, year)

	assert.Equal(t, firstDate.GetFormattedDate(), expectedFormattedDate)
}

func TestFirstDateOfTheMouthWhenTheDataAreInvalid(t *testing.T) {
	month := Month{Content: 13}
	year := NewYear(2018)

	_, err := FirstDateOfMonth(month, year)
	assert.NotEqual(t, nil, err)
}

func TestLastDateOfTheMouth(t *testing.T) {
	expectedFormattedDate := "30/04/2018"

	month := APRIL()
	year := NewYear(2018)			
	lastDate, _ := LastDateOfMonth(month, year)
	
	assert.Equal(t, lastDate.GetFormattedDate(), expectedFormattedDate)
}


func TestLastDateOfTheMouthWhenTheDataAreInvalid(t *testing.T) {
	month := Month{Content: 13}
	year := NewYear(2018)

	_, err := LastDateOfMonth(month, year)
	assert.NotEqual(t, nil, err)
}

// DateOf truncates to the calendar day in UTC, whatever the input's zone —
// the single notion of "today" shared by pause/resume stamping and the daily
// generation engine (ADR 0005).
func TestDateOfTruncatesToTheUTCCalendarDay(t *testing.T) {
	cest := time.FixedZone("CEST", 2*60*60)

	d := DateOf(time.Date(2026, 9, 25, 0, 30, 0, 0, cest))

	assert.Equal(t, "2026-09-24", d.GetIsoFormattedDate())
	assert.Equal(t, time.Date(2026, 9, 24, 0, 0, 0, 0, time.UTC), d.GetTime())
}
