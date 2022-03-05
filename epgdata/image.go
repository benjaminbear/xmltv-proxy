package epgdata

import (
	"encoding/json"
)

type Images struct {
	Size130 *Image
	Size320 *Image
	Size476 *Image
	Size952 *Image
}

type Image struct {
	Source string
	Width  string
	Height string
}

func (i *Images) UnmarshalJSON(data []byte) error {
	var aux []struct {
		Size1 string `json:"size1"`
		Size2 string `json:"size2"`
		Size3 string `json:"size3"`
		Size4 string `json:"size4"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if len(aux) > 0 {
		i.Size130 = &Image{
			Source: aux[0].Size1,
			Width:  "130",
			Height: "101",
		}

		i.Size320 = &Image{
			Source: aux[0].Size2,
			Width:  "320",
			Height: "250",
		}

		i.Size476 = &Image{
			Source: aux[0].Size3,
			Width:  "476",
			Height: "357",
		}

		i.Size952 = &Image{
			Source: aux[0].Size4,
			Width:  "952",
			Height: "714",
		}
	}

	return nil
}
