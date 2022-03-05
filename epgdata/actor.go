package epgdata

import (
	"encoding/json"
)

type Actor struct {
	Actor string
	Role  string
}

func (a *Actor) UnmarshalJSON(data []byte) error {
	var aux map[string]string

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	for role, actor := range aux {
		a.Role = role
		a.Actor = actor
	}

	return nil
}
