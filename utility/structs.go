package utility

import (
	"encoding/json"

	"github.com/fatih/structs"
)

type Structs struct{}

func (s *Structs) ToMap(v interface{}) map[string]interface{} {
	m := structs.Map(v)
	return m
}

func (s *Structs) Attr(v interface{}) map[string]string {
	y := s.ToMap(v)
	m := make(map[string]string)
	for k, val := range y {
		j, _ := json.MarshalIndent(val, "", "  ")
		m[k] = string(j)
	}

	attrs := Attr(v)

	for k, v := range attrs {
		m[k] = v
	}

	return m
}
