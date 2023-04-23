package config

import "reflect"

type Field struct {
	Name        string
	Description string
	Style       string
	Tag         string
	Type        reflect.Type
}
