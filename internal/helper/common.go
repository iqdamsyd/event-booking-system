package helper

import "time"

func ParseYYYYMMDD(stringDate string) (*time.Time, error) {
	layout := "2006-01-02"
	parsedTime, err := time.Parse(layout, stringDate)
	if err != nil {
		return nil, err
	}

	return &parsedTime, nil
}

func ParseYYYYMMDD2359(stringDate string) (*time.Time, error) {
	layout := "2006-01-02"
	parsedTime, err := time.Parse(layout, stringDate)
	if err != nil {
		return nil, err
	}

	endOfDay := time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(), 23, 59, 59, 0, parsedTime.Location())

	return &endOfDay, nil
}
