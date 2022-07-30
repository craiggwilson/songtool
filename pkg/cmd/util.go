package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/craiggwilson/songtool/pkg/theory/key"
)

func firstNonEmptyString(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}

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

type keySurrogate struct {
	Name        string   `json:"name"`
	DegreeClass int      `json:"degreeClass"`
	PitchClass  int      `json:"pitchClass"`
	Kind        key.Kind `json:"kind"`
}

type noteSurrogate struct {
	Name        string `json:"name"`
	Interval    string `json:"interval"`
	DegreeClass int    `json:"degreeClass"`
	PitchClass  int    `json:"pitchClass"`
}

type scaleSurrogate struct {
	Name  string          `json:"name"`
	Notes []noteSurrogate `json:"notes"`
}
