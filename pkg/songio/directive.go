package songio

import (
	"encoding/json"

	"github.com/craiggwilson/songtool/pkg/theory/key"
)

type KeyDirectiveLine struct {
	Key key.Parsed `json:"key"`
}

func (d *KeyDirectiveLine) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Directive string     `json:"directive"`
		Value     key.Parsed `json:"value"`
	}{
		Directive: "key",
		Value:     d.Key,
	})
}

func (d *KeyDirectiveLine) line() {}

type SectionEndDirectiveLine struct {
	Name string `json:"name"`
}

func (d *SectionEndDirectiveLine) MarshalJSON() ([]byte, error) {
	return json.Marshal(UnknownDirectiveLine{
		Name:  "sectionEnd",
		Value: d.Name,
	})
}

func (d *SectionEndDirectiveLine) line() {}

type SectionStartDirectiveLine struct {
	Name string `json:"name"`
}

func (d *SectionStartDirectiveLine) MarshalJSON() ([]byte, error) {
	return json.Marshal(UnknownDirectiveLine{
		Name:  "sectionStart",
		Value: d.Name,
	})
}

func (d *SectionStartDirectiveLine) line() {}

type TitleDirectiveLine struct {
	Title string `json:"title"`
}

func (d *TitleDirectiveLine) MarshalJSON() ([]byte, error) {
	return json.Marshal(UnknownDirectiveLine{
		Name:  "title",
		Value: d.Title,
	})
}

func (d *TitleDirectiveLine) line() {}

type UnknownDirectiveLine struct {
	Name  string `json:"directive"`
	Value string `json:"value,omitempty"`
}

func (d *UnknownDirectiveLine) line() {}
