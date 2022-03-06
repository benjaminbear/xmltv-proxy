package epgdata

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type seriesNumber struct {
	Valid  bool
	Number int
	Total  int
}

func (s *seriesNumber) UnmarshalJSON(data []byte) error {
	var aux string

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	reg := regexp.MustCompile("[^0-9/]+")
	aux = reg.ReplaceAllString(aux, "")

	if aux == "0" || aux == "" {
		return nil
	}

	// is number e.g. "1"
	s.Number, err = strconv.Atoi(aux)
	if err == nil {
		s.Valid = true

		return nil
	}

	// is format "1/3"
	if strings.Contains(aux, "/") {
		splits := strings.Split(aux, "/")

		if len(splits) == 2 {
			s.Number, err = strconv.Atoi(splits[0])
			if err != nil {
				return err
			}

			s.Total, err = strconv.Atoi(splits[1])
			if err != nil {
				return err
			}

			s.Valid = true

			return nil
		}
	}

	return nil
}

// XMLTVString returns a String prepared for XMLTv and based to 0.
func (s *seriesNumber) XMLTVString() string {
	if !s.Valid {
		// s.Number == 0
		return fmt.Sprintf("%d", s.Number)
	}

	if s.Total > 0 {
		return fmt.Sprintf("%d/%d", s.Number-1, s.Total-1)
	}

	return fmt.Sprintf("%d", s.Number-1)
}
