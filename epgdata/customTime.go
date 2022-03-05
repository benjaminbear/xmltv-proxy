package epgdata

import (
	"encoding/json"
	"time"
)

type dateTime struct {
	time.Time
}

func (t *dateTime) UnmarshalJSON(data []byte) error {
	var aux int64

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Unix timestamp has a fixed timezone
	*t = dateTime{time.Unix(aux, 0)}

	return nil
}

/*
type tinyDate struct {
	time.Time
}

func (t *tinyDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const shortForm = "20060102"
	var v string

	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}

	parse, err := time.Parse(shortForm, v)
	if err != nil {
		return err
	}

	*t = tinyDate{parse}

	return nil
}
*/
