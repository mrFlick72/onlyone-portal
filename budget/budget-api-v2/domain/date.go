package domain

type Date struct {
	content string
}

func (d *Date) GetFormattedDate() (string, error) {
	return "", nil
}

func (d *Date) GetIsoFormattedDate() (string, error) {
	return "", nil
}

func firstDateOfMonth(month Month, year Year) (*Date, error) {
	return nil, nil
}

func lastDateOfMonth(month Month, year Year) (*Date, error) {
	return nil, nil
}
