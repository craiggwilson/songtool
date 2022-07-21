package songio

import "github.com/craiggwilson/songtools/theory"

type KeyDirectiveLine struct {
	Key theory.Key
}

func (d *KeyDirectiveLine) line() {}

type SectionEndDirectiveLine struct {
	Name string
}

func (d *SectionEndDirectiveLine) line() {}

type SectionStartDirectiveLine struct {
	Name string
}

func (d *SectionStartDirectiveLine) line() {}

type TitleDirectiveLine struct {
	Title string
}

func (d *TitleDirectiveLine) line() {}

type UnknownDirectiveLine struct {
	Name  string
	Value string
}

func (d *UnknownDirectiveLine) line() {}
