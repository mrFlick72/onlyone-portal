package revenue

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
)

func TestParseYearQueryParam(t *testing.T) {
	year, err := parseYearQueryParam("year=2023")
	assert.Equal(t, nil, err)
	assert.Equal(t, date.NewYear(2023), year)
}

func TestParseYearQueryParamWithTrailingSegments(t *testing.T) {
	year, err := parseYearQueryParam("year=2023;other=foo")
	assert.Equal(t, nil, err)
	assert.Equal(t, date.NewYear(2023), year)
}

func TestParseYearQueryParamEmpty(t *testing.T) {
	_, err := parseYearQueryParam("")
	assert.NotEqual(t, nil, err)
}

func TestParseYearQueryParamMissingEquals(t *testing.T) {
	_, err := parseYearQueryParam("year")
	assert.NotEqual(t, nil, err)
}

func TestParseYearQueryParamWrongKey(t *testing.T) {
	_, err := parseYearQueryParam("month=2023")
	assert.NotEqual(t, nil, err)
}
