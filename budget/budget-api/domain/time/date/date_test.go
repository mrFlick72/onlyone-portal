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
// scheduled expense job (ADR 0005).
func TestDateOfTruncatesToTheUTCCalendarDay(t *testing.T) {
	cest := time.FixedZone("CEST", 2*60*60)

	d := DateOf(time.Date(2026, 9, 25, 0, 30, 0, 0, cest))

	assert.Equal(t, "2026-09-24", d.GetIsoFormattedDate())
	assert.Equal(t, time.Date(2026, 9, 24, 0, 0, 0, 0, time.UTC), d.GetTime())
}

func TestAddDaysCrossesMonthAndYearBoundaries(t *testing.T) {
	d, _ := IsoDateFor("2026-12-31")

	next := d.AddDays(1)

	assert.Equal(t, "2027-01-01", next.GetIsoFormattedDate())
	assert.Equal(t, "2026-12-31", d.GetIsoFormattedDate())
}

func TestDaysInMonth(t *testing.T) {
	for iso, expected := range map[string]int{
		"2026-04-10": 30,
		"2026-02-10": 28,
		"2028-02-10": 29,
		"2026-01-10": 31,
	} {
		d, _ := IsoDateFor(iso)
		assert.Equal(t, expected, d.DaysInMonth())
	}
}

func TestIsAfter(t *testing.T) {
	d, _ := IsoDateFor("2026-09-24")
	next := d.AddDays(1)

	assert.Equal(t, true, next.IsAfter(*d))
	assert.Equal(t, false, d.IsAfter(next))
	assert.Equal(t, false, d.IsAfter(*d))
}
