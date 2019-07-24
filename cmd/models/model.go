package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/flect"
)

// Model represents a resource model object
type Model struct {
	Name          string               `json:"name"`
	Type          string               `json:"type"`
	Ident         flect.Ident          `json:"ident"`
	Attributes    map[string]Attribute `json:"attributes"`
	Nested        []Model              `json:"nested"`
	Imports       []string             `json:"imports"`
	KeySchema     map[string]string    `json:"key_schema"`
	GeneratedID   bool                 `json:"generated_id"`
	CompositeKey  bool                 `json:"composite_key"`
	BillingMode   string               `json:"billing_mode"`
	CapacityUnits map[string]int64     `json:"capacity_units"`
}

// Attribute represents a resource model's attribute
type Attribute struct {
	Name    string      `json:"name"`
	Ident   flect.Ident `json:"ident"`
	GoType  string      `json:"go_type"`
	AwsType string      `json:"aws_type"`
}

// New returns a new model object
func New(name string, slice bool, attributes string, options map[string]interface{}) Model {
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

	// parse nested models
	attributes = m.parseNested(attributes)
	m.parseAttributes(attributes)

	// handle all option values
	var id, withDates, softDelete bool
	var keySchema, billing string
	var capacity map[string]int64
	if options != nil {
		id = options["id"].(bool)
		withDates = options["dates"].(bool)
		softDelete = options["softDelete"].(bool)
		keySchema = options["keySchema"].(string)
		billing = options["billing"].(string)
		capacity = options["capacity"].(map[string]int64)
	} else {
		id, withDates, softDelete = false, false, false
		billing = "provisioned"
		capacity = map[string]int64{
			"read":  1,
			"write": 1,
		}
	}

	if id {
		m.GeneratedID = true
		a := Attribute{Name: "id", Ident: flect.New("id"), AwsType: "S", GoType: "string"}
		m.Imports = appendStringIfMissing(m.Imports, "github.com/gofrs/uuid")
		m.addAttribute(a)
		m.KeySchema = map[string]string{
			"HASH": "id",
		}
	} else if len(keySchema) > 0 {
		m.parseKeySchema(keySchema)
	}

	if withDates {
		m.Imports = appendStringIfMissing(m.Imports, "time")
		m.addAttribute(Attribute{Name: "createdAt", Ident: flect.New("createdAt"), AwsType: "S", GoType: "time.Time"})
		m.addAttribute(Attribute{Name: "updatedAt", Ident: flect.New("updatedAt"), AwsType: "S", GoType: "time.Time"})
	}

	if softDelete {
		m.Imports = appendStringIfMissing(m.Imports, "time")
		m.addAttribute(Attribute{Name: "deletedAt", Ident: flect.New("deletedAt"), AwsType: "S", GoType: "*time.Time"})
	}

	m.BillingMode = strings.ToLower(billing)

	if m.BillingMode == "provisioned" {
		m.CapacityUnits = capacity
	}

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
		clAttr = strings.Replace(clAttr, np, "", 1)
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
	nmn := strings.Replace(strings.TrimSpace(attributes[cI:bI-1]), ",", "", 1)
	attr := attributes[bI+1 : pos]
	n := New(nmn, slice, attr, nil)

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
		)

		if len(inputs) > 1 {
			goType = inputs[1]
		}

		attr := Attribute{
			Name:    name,
			Ident:   flect.New(name),
			GoType:  goType,
			AwsType: awsType(goType),
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
		m.Imports = appendStringIfMissing(m.Imports, "github.com/gofrs/uuid")
	}
}

// GetImports recursively iterates through all import slices and adds the import to the root model
func (m *Model) GetImports() []string {
	var imports []string
	if len(m.Nested) > 0 {
		for _, n := range m.Nested {
			// get all imports of the nested model
			nI := n.GetImports()

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
		if m.Attributes == nil {
			m.Attributes = map[string]Attribute{
				a.Name: a,
			}
		}
		m.Attributes[a.Name] = a
	}

}

// parseKeySchema parses a given keySchema and add it to the model
func (m *Model) parseKeySchema(schema string) {
	for _, k := range strings.Split(schema, ",") {
		key := strings.Split(k, ":")
		if m.KeySchema == nil {
			m.KeySchema = map[string]string{
				strings.ToUpper(key[1]): key[0],
			}
		} else {
			m.KeySchema[strings.ToUpper(key[1])] = key[0]
		}
	}

	if c, err := m.checkKeys(); !c {
		log.Fatal(err)
	}

}

// checkKeys checks the Key Schema of the model against its attributes
func (m *Model) checkKeys() (bool, error) {
	check := map[string]byte{
		"hash":  0,
		"range": 0,
	}

	hashKey := m.KeySchema["HASH"]
	rangeKey := m.KeySchema["RANGE"]

	for _, a := range m.Attributes {
		if a.Name == hashKey {
			check["hash"]++
		}
		if a.Name == rangeKey {
			check["range"]++
		}
	}

	if check["hash"] == 0 {
		return false, fmt.Errorf("No Hash Key defined for %s. Cannot identify ID Attribute", m.Name)
	}

	if check["hash"] == 1 {
		if check["range"] >= 1 {
			m.CompositeKey = true
		}

		return true, nil
	}

	return false, fmt.Errorf("Too many keys defined in Key Schema")

}

// GetConfigs returns the MUGConfig and ServerlessConfig for this Model
func (m Model) GetConfigs() (MUGConfig, ServerlessConfig) {
	attributeDefinitions := map[string]AttributeDefinition{}
	for _, k := range m.KeySchema {
		a := m.Attributes[k]
		if len(a.Name) > 0 {
			attributeDefinitions[a.Name] = AttributeDefinition{
				Ident:   a.Ident,
				AwsType: a.AwsType,
			}
		}
	}

	// update mug.config.json
	r := &NewResource{
		Ident:      flect.New(m.Name),
		Attributes: attributeDefinitions,
	}
	mc := ReadMUGConfig()
	mc.Resources[m.Name] = r
	// mc.Write()

	// update serverless.yml
	var path string
	singular := m.Ident.Singularize().String()
	plural := m.Ident.Pluralize().String()
	if m.CompositeKey {
		path = fmt.Sprintf("%s/{%s}/{%s}", plural, m.KeySchema["HASH"], m.KeySchema["RANGE"])
	} else {
		path = fmt.Sprintf("%s/{%s}", plural, m.KeySchema["HASH"])
	}
	fns := []*Function{
		&Function{Name: "create" + "_" + singular, Handler: "create", Path: plural, Method: "post"},
		&Function{Name: "read" + "_" + singular, Handler: "read", Path: path, Method: "get"},
		&Function{Name: "update" + "_" + singular, Handler: "update", Path: path, Method: "put"},
		&Function{Name: "delete" + "_" + singular, Handler: "delete", Path: path, Method: "delete"},
		&Function{Name: "list" + "_" + plural, Handler: "list", Path: plural, Method: "get"},
	}

	sc := mc.NewServerlessConfig()
	sc.SetResourceWithModel(r, m)
	sc.SetFunctions(fns)

	return mc, sc
}

// Write write the Model definition to the modelName.json
func (m Model) Write(path string) {
	json, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(filepath.Join(path, "functions", m.Name, fmt.Sprintf("%s.json", m.Name)), json, 0644)
	if err != nil {
		log.Fatal(err)
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
