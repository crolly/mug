package models

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/gobuffalo/flect"

	"github.com/crolly/oas3"
)

// OAS3Spec is the extenstion of oas3.Document
type OAS3Spec struct {
	*oas3.Document
}

// GetOAS returns the projects oas3.Document or a new one with default values
func GetOAS(path string) *OAS3Spec {
	s, err := oas3.LoadFile(filepath.Join(path, "spec.yml"))
	if err != nil && s != nil {
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

func (s *OAS3Spec) addPath(f *ServerlessFunction, i flect.Ident) {
	e := f.Events[0].HTTP

	if e != nil {
		s.initPaths(e)
		switch strings.TrimPrefix(f.Handler, "bin/") {
		case "create":
			s.addCreatePath(f, i)
		case "read":
			s.addReadPath(f, i)
		case "update":
			s.addUpdatePath(f, i)
		case "delete":
			s.addDeletePath(f, i)
		case "list":
			s.addListPath(f, i)
		}
	}
}

func (s *OAS3Spec) initPaths(e *HTTPEvent) {
	if len(s.Document.Paths) == 0 {
		s.Document.Paths = oas3.Paths{}
	}

	if s.Document.Paths["/"+e.Path] == nil {
		s.Document.Paths["/"+e.Path] = &oas3.PathItem{}
	}
}

func (s *OAS3Spec) addCreatePath(f *ServerlessFunction, i flect.Ident) {
	e := f.Events[0].HTTP
	singular := i.Singularize().Pascalize().String()
	s.Document.Paths["/"+e.Path].Post = &oas3.Operation{
		Summary:     "Create " + singular,
		Description: fmt.Sprintf("Creates a new %s object.", singular),
		OperationID: f.Name,
		Tags:        []string{i.Pluralize().Camelize().String()},
		RequestBody: defaultRequestBody(singular),
		Responses: map[string]*oas3.Response{
			"200": default200Response("JSON Object of the ", singular),
		},
	}
}

func (s *OAS3Spec) addReadPath(f *ServerlessFunction, i flect.Ident) {
	e := f.Events[0].HTTP
	singular := i.Singularize().Pascalize().String()
	s.Document.Paths["/"+e.Path].Get = &oas3.Operation{
		Summary:     "Read " + singular,
		Description: fmt.Sprintf("Retrieves the details of an existing %s. Supply the identifier from either the creation request or the %s list.", singular, singular),
		OperationID: f.Name,
		Tags:        []string{i.Pluralize().Camelize().String()},
		Parameters: []*oas3.Parameter{
			&oas3.Parameter{
				Name:        "id",
				In:          "path",
				Required:    true,
				Description: "The ID of the " + singular + " to retrieve",
				Schema: &oas3.Schema{
					Type: "string",
				},
			},
		},
		Responses: map[string]*oas3.Response{
			"200": default200Response("JSON Object of the ", singular),
			"404": &oas3.Response{
				Description: "Response if the " + singular + " with the given ID could not be found.",
			},
		},
	}
}

func (s *OAS3Spec) addUpdatePath(f *ServerlessFunction, i flect.Ident) {
	e := f.Events[0].HTTP
	singular := i.Singularize().Pascalize().String()
	s.Document.Paths["/"+e.Path].Put = &oas3.Operation{
		Summary:     "Update " + singular,
		Description: fmt.Sprintf("Updates the specific %s by setting the values of the parameters passed.", singular),
		OperationID: f.Name,
		Tags:        []string{i.Pluralize().Camelize().String()},
		Parameters: []*oas3.Parameter{
			&oas3.Parameter{
				Name:        "id",
				In:          "path",
				Required:    true,
				Description: "The ID of the " + singular + " to update",
				Schema: &oas3.Schema{
					Type: "string",
				},
			},
		},
		RequestBody: defaultRequestBody(singular),
		Responses: map[string]*oas3.Response{
			"200": default200Response("JSON Object of the ", singular),
		},
	}
}

func (s *OAS3Spec) addDeletePath(f *ServerlessFunction, i flect.Ident) {
	e := f.Events[0].HTTP
	singular := i.Singularize().Pascalize().String()
	s.Document.Paths["/"+e.Path].Delete = &oas3.Operation{
		Summary:     "Delete " + singular,
		Description: fmt.Sprintf("Delete a %s.", singular),
		OperationID: f.Name,
		Tags:        []string{i.Pluralize().Camelize().String()},
		Parameters: []*oas3.Parameter{
			&oas3.Parameter{
				Name:        "id",
				In:          "path",
				Required:    true,
				Description: "The ID of the " + singular + " to delete",
				Schema: &oas3.Schema{
					Type: "string",
				},
			},
		},
		Responses: map[string]*oas3.Response{
			"200": default200Response("Message: Deleted ", singular),
		},
	}
}

func (s *OAS3Spec) addListPath(f *ServerlessFunction, i flect.Ident) {
	e := f.Events[0].HTTP
	plural := i.Pluralize().Pascalize().String()
	s.Document.Paths["/"+e.Path].Get = &oas3.Operation{
		Summary:     "List " + plural,
		Description: fmt.Sprintf("Returns a list of %s. The %s are returned unsorted.", plural, plural),
		OperationID: f.Name,
		Tags:        []string{i.Pluralize().Camelize().String()},
		Responses: map[string]*oas3.Response{
			"200": default200Response("Array of ", plural),
		},
	}
}

func default200Response(desc, term string) *oas3.Response {
	return &oas3.Response{
		Description: desc + term,
		Content: map[string]*oas3.MediaType{
			"application/json": &oas3.MediaType{
				Schema: &oas3.Schema{
					Ref: "#/components/schemas/" + term,
				},
			},
		},
	}
}

func defaultRequestBody(term string) *oas3.RequestBody {
	return &oas3.RequestBody{
		Required: true,
		Content: map[string]*oas3.MediaType{
			"application/json": &oas3.MediaType{
				Schema: &oas3.Schema{
					Ref: "#/components/schemas/" + term,
				},
			},
		},
	}
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
