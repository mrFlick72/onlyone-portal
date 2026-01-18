package date

type Month struct {
	Content int8
}

func JANUARY() Month {
	return Month{
		Content: 1,
	}
}

func FEBRUARY() Month {
	return Month{
		Content: 2,
	}
}

func MARCH() Month {
	return Month{
		Content: 3,
	}
}

func APRIL() Month {
	return Month{
		Content: 4,
	}
}

func MAY() Month {
	return Month{
		Content: 5,
	}
}

func JUNE() Month {
	return Month{
		Content: 6,
	}
}

func JULY() Month {
	return Month{
		Content: 7,
	}
}

func AUGUST() Month {
	return Month{
		Content: 8,
	}
}

func SEPTEMBER() Month {
	return Month{
		Content: 9,
	}
}

func OCTOBER() Month {
	return Month{
		Content: 10,
	}
}

func NOVEMBER() Month {
	return Month{
		Content: 11,
	}
}

func DECEMBER() Month {
	return Month{
		Content: 12,
	}
}

type Year struct {
	Content int16
}

func NewYear(year int16) Year {
	return Year{
		Content: year,
	}
}
