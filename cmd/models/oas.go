package models

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/crolly/oas3"
	"github.com/ghodss/yaml"
)

// OAS3Spec is the extenstion of oas3.Document
type OAS3Spec struct {
	*oas3.Document
}

// GetOAS returns the projects oas3.Document or a new one with default values
func GetOAS(path string) *OAS3Spec {
	s, err := oas3.LoadFile(filepath.Join(path, "spec.yml"))
	if err != nil {
		return &OAS3Spec{
			Document: s,
		}
	}

	return &OAS3Spec{
		Document: &oas3.Document{
			Version: "3.0.0",
			Info: &oas3.Info{
				Version: "1.0.0",
			},
		},
	}
}

func (s *OAS3Spec) Write(path string) {
	yml, err := yaml.Marshal(s.Document)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(filepath.Join(path, "spec.yml"), yml, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *OAS3Spec) addPath(f ServerlessFunction) {

}

func (s *OAS3Spec) addComponent(m Model) {
	req := []string{}
	for _, k := range m.KeySchema {
		req = append(req, k)
	}

	props := map[string]*oas3.Schema{}
	for n, a := range m.Attributes {
		props[n] = getPropDef(a.GoType)
	}

	singular := m.Ident.Singularize().Pascalize().String()
	plural := m.Ident.Pluralize().Pascalize().String()

	if s.Document.Components == nil {
		s.Document.Components = &oas3.Components{
			Schemas: map[string]*oas3.Schema{},
		}
	}
	s.Document.Components.Schemas[singular] = &oas3.Schema{
		Required:   req,
		Properties: props,
	}
	s.Document.Components.Schemas[plural] = &oas3.Schema{
		Type: "array",
		Items: &oas3.Schema{
			Ref: "#/components/schemas/" + singular,
		},
	}
}

func (s *OAS3Spec) setTitle(title string) {
	s.Document.Info.Title = title
}

func getPropDef(attrType string) *oas3.Schema {
	s := &oas3.Schema{}

	var f string
	var t string
	switch attrType {
	case "*int32", "*int64":
		t = "integer"
	case "*float32":
		f = "float"
		t = "number"
	case "*float64":
		f = "double"
		t = "number"
	case "*byte":
		f = "byte"
		t = "string"
	case "*bool":
		t = "boolean"
	case "time.Time", "*time.Time":
		f = "date-time"
		t = "string"
	default:
		t = attrType
	}

	if len(f) > 0 {
		s.Format = f
	}
	s.Type = t

	return s
}
