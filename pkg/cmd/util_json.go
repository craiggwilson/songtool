package cmd

import (
	"encoding/json"
	"fmt"
)

func marshalJSON(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", " ")
}

func printJSON(v interface{}) error {
	out, err := marshalJSON(v)
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}

type scaleSurrogate struct {
	Name  string          `json:"name"`
	Notes []noteSurrogate `json:"notes"`
}

type noteSurrogate struct {
	Name        string `json:"name"`
	Interval    string `json:"interval"`
	DegreeClass int    `json:"degreeClass"`
	PitchClass  int    `json:"pitchClass"`
}
