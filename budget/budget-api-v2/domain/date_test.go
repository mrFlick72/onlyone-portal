package domain

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

//     @Test
//     public void dateFromString() {
//         Date expectedDateForDateString = new Date(LocalDate.of(2018, 02, 25));
//         Date actualDateForDateString = Date.dateFor("25/02/2018");
//         Assertions.assertEquals(actualDateForDateString, expectedDateForDateString);
//     }

//     @Test
//     public void firstDateOfMonth() {
//         Date expected = Date.dateFor("01/02/2018");
//         Date actual = Date.firstDateOfMonth(Month.FEBRUARY, Year.of(2018));

//         Assertions.assertEquals(actual, expected);
//     }

//     @Test
//     public void lastDateOfMonth() {
//         Date expected = Date.dateFor("28/02/2018");
//         Date actual = Date.lastDateOfMonth(Month.FEBRUARY, Year.of(2018));

//         Assertions.assertEquals(actual, expected);
//     }
// }

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
