package epgdata

import (
	"encoding/json"
	"strings"
)

type StdListElement struct {
	List []string
}

func (r *StdListElement) UnmarshalJSON(data []byte) error {
	var aux string
	var aux2 []string

	err := json.Unmarshal(data, &aux)
	if err == nil {
		r.splitElements(aux)

		return nil
	}

	err = json.Unmarshal(data, &aux2)
	if err != nil {
		return err
	}

	r.splitElements(aux2[0])

	return nil
}

func (r *StdListElement) splitElements(v string) {
	splits := strings.Split(v, ",")

	for _, part := range splits {
		part = strings.TrimSpace(part)
		r.List = append(r.List, part)
	}
}
