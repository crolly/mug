package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gobuffalo/flect"
)

// ResourceConfig represents mu's configuration for resources
type ResourceConfig struct {
	ProjectName string              `json:"projectName"`
	ProjectPath string              `json:"projectPath"`
	ImportPath  string              `json:"importPath"`
	Region      string              `json:"region"`
	Resources   map[string]Resource `json:"resources"`
}

// Resource represents a single Resource of the project's config
type Resource struct {
	Ident     flect.Ident `json:"ident"`
	Functions []Function  `json:"functions"`
}

// Function represents a Function
type Function struct {
	Name    string `json:"name"`
	Handler string `json:"handler"`
	Path    string `json:"path"`
	Method  string `json:"method"`
}

func readConfig() ResourceConfig {
	wd := getWorkingDir()

	configFile, err := os.Open(filepath.Join(wd, "mug.config.json"))
	if err != nil {
		log.Fatal(err)
	}

	defer configFile.Close()

	data, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatal(err)
	}

	var config ResourceConfig

	json.Unmarshal(data, &config)

	// make sure map exists
	if len(config.Resources) == 0 {
		config.Resources = make(map[string]Resource)
	}

	return config
}

func writeConfig(config ResourceConfig) {
	// wd := getWorkingDir()

	fileName := filepath.Join(config.ProjectPath, "mug.config.json")

	configJSON, _ := json.MarshalIndent(config, "", "  ")
	_ = ioutil.WriteFile(fileName, configJSON, 0644)
}

// Model represents a resource model object
type Model struct {
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Ident      flect.Ident `json:"ident"`
	Attributes []Attribute `json:"attributes"`
	Nested     []Model     `json:"nested"`
	Imports    []string    `json:"imports"`
}

// Attribute represents a resource model's attribute
type Attribute struct {
	Name    string      `json:"name"`
	Ident   flect.Ident `json:"ident"`
	GoType  string      `json:"goType"`
	AwsType string      `json:"awsType"`
	Hash    bool        `json:"hash"`
}

// returns a new model object
func newModel(name string, slice bool, attributes string, withID bool, withDates bool) Model {
	ident := flect.New(name)
	m := Model{
		Name:  ident.Camelize().String(),
		Ident: ident,
	}

	if slice {
		m.Type = fmt.Sprintf("[]%s", m.Ident.Pascalize())
	} else {
		m.Type = m.Ident.Pascalize().String()
	}

	if withID {
		m.Imports = appendStringIfMissing(m.Imports, "github.com/gofrs/uuid")
		m.addAttribute(Attribute{Name: "id", Ident: flect.New("id"), AwsType: "S", GoType: "uuid.UUID", Hash: true})
	}

	if withDates {
		m.Imports = appendStringIfMissing(m.Imports, "time")
		m.addAttribute(Attribute{Name: "createdAt", Ident: flect.New("createdAt"), AwsType: "S", GoType: "time.Time", Hash: false})
		m.addAttribute(Attribute{Name: "updatedAt", Ident: flect.New("updatedAt"), AwsType: "S", GoType: "time.Time", Hash: false})
	}

	// parse nested models
	attributes = m.parseNested(attributes)
	m.parseAttributes(attributes)

	return m
}

// parseNested parses the attributes string for nested models
func (m *Model) parseNested(attributes string) string {
	var (
		cob    []int        // curly opening bracket slice to remember position
		cbc    = 0          // closing curly bracket counter
		sob    []int        // square opening bracket slice to remember position
		sbc    = 0          // closing square bracket counter
		rm     []string     // string slice with nested parts to remove
		clAttr = attributes // cleared attribute string without nested parts
	)
	for pos, char := range attributes {
		if char == '{' {
			// opening bracket
			cob = append(cob, pos)
		}
		if char == '}' {
			// closing bracket
			cbc++
		}
		if char == '[' {
			sob = append(sob, pos)
		}
		if char == ']' {
			sbc++
		}

		if len(cob) > 0 && len(cob) == cbc { // found single nested
			cI := m.addNested(cob, pos, attributes, false)

			// append nested part to rm slice
			rm = append(rm, attributes[cI:pos+1])

			cob = nil
			cbc = 0
		}

		if len(sob) > 0 && len(sob) == sbc { // found slice nested
			cI := m.addNested(sob, pos, attributes, true)

			// append nested part to rm slice
			rm = append(rm, attributes[cI:pos+1])

			sob = nil
			sbc = 0
		}
	}

	for _, np := range rm {
		clAttr = strings.ReplaceAll(clAttr, np, "")
	}

	return clAttr
}

// addNested adds a nested model to the resource model
func (m *Model) addNested(b []int, pos int, attributes string, slice bool) int {
	// opening bracket index
	bI := b[0]
	// comma index
	cI := strings.LastIndex(attributes[0:bI-1], ",")
	if cI < 0 {
		cI = 0
	}

	// new model name ensured to not have a comma or spaces
	nmn := strings.ReplaceAll(strings.TrimSpace(attributes[cI:bI-1]), ",", "")
	attr := attributes[bI+1 : pos]
	n := newModel(nmn, slice, attr, false, false)

	m.Nested = append(m.Nested, n)

	return cI
}

// parseAttributes parses all the attributes attached to a resource model
func (m *Model) parseAttributes(attrs string) {
	for _, a := range strings.Split(attrs, ",") {
		inputs := strings.Split(a, ":")
		fmt.Println(inputs)
		name := inputs[0]

		// handle optional inputs
		var (
			goType = "string"
			hash   = false
			err    error
		)

		if len(inputs) > 1 {
			goType = inputs[1]
			if len(inputs) > 2 {
				hash, err = strconv.ParseBool(inputs[2])
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		attr := Attribute{
			Name:    name,
			Ident:   flect.New(name),
			GoType:  goType,
			AwsType: awsType(goType),
			Hash:    hash,
		}

		m.addImport(goType)

		m.addAttribute(attr)
	}
}

// addImport will add an import directive if the given type requires it
func (m *Model) addImport(goType string) {
	switch goType {
	case "time.Time", "*time.Time":
		m.Imports = appendStringIfMissing(m.Imports, "time")
	case "uuid.UUID":
		m.Imports = appendStringIfMissing(m.Imports, "github.com/satori/go.uuid")
	}
}

// getImports recursively iterates through all import slices and adds the import to the root model
func (m *Model) getImports() []string {
	var imports []string
	if len(m.Nested) > 0 {
		for _, n := range m.Nested {
			// get all imports of the nested model
			nI := n.getImports()

			// iterate over imports and append new ones to imports slice
			for _, i := range nI {
				imports = appendStringIfMissing(imports, i)
			}
		}
	}

	for _, i := range m.Imports {
		imports = appendStringIfMissing(imports, i)
	}

	return imports
}

// addAttribute adds an attribute to a resource model
func (m *Model) addAttribute(a Attribute) {
	// make sure all attributes have names
	if a.Name != "" {
		m.Attributes = append(m.Attributes, a)
	}

}

// String prints a representation of a model
func (m Model) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("// %s defines the %s model\n", m.Ident.Pascalize(), m.Ident.Pascalize()))
	sb.WriteString(fmt.Sprintf("type %s struct {\n", m.Ident.Pascalize()))
	for _, a := range m.Attributes {
		sb.WriteString(fmt.Sprintf("%s\n", a.String()))
	}
	if len(m.Nested) > 0 {
		sb.WriteString("\n")
		for _, n := range m.Nested {
			sb.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\"`\n", n.Ident.Pascalize(), n.Type, n.Ident.Underscore()))
		}
		sb.WriteString("}\n")
		sb.WriteString("\n")
		for _, n := range m.Nested {
			sb.WriteString(n.String())
			sb.WriteString("\n")
		}

	} else {
		sb.WriteString("}\n")
	}

	return sb.String()
}

// String returns the string representation of an attribute
func (a Attribute) String() string {
	return fmt.Sprintf("\t%s %s `json:\"%s\"`", a.Ident.Pascalize(), a.GoType, a.Ident.Underscore())
}
