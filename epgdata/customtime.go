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
