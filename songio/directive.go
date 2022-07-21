package songio

import "github.com/craiggwilson/songtools/theory"

type Directive interface {
	dir()
}

type KeyDirective struct {
	Key theory.Key
}

func (d *KeyDirective) Name() string { return "key" }
func (d *KeyDirective) dir()         {}

type SectionEndDirective struct {
	SectionName string
}

func (d *SectionEndDirective) Name() string { return d.SectionName }
func (d *SectionEndDirective) dir()         {}

type SectionStartDirective struct {
	SectionName string
}

func (d *SectionStartDirective) Name() string { return d.SectionName }
func (d *SectionStartDirective) dir()         {}

type TitleDirective struct {
	Title string
}

func (d *TitleDirective) Name() string { return "title" }
func (d *TitleDirective) dir()         {}

type UnknownDirective struct {
	DirectiveName string
	Value         string
}

func (d *UnknownDirective) Name() string { return d.DirectiveName }
func (d *UnknownDirective) dir()         {}
