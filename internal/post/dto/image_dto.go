package dto

import (
	"encoding/json"
	"io"
)

type Image struct {
	FileName string `json:"fileName"`
	Data     string `json:"file"`
}

func (i *Image) FromJSON(r io.Reader) error {
	err := json.NewDecoder(r).Decode(i)
	return err
}
