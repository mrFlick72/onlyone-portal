package date

import (
	"testing"

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
