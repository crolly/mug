package models

// TemplateConfig ...
type TemplateConfig struct {
	Transform string                 `yaml:"Transform"`
	Resources map[string]SAMResource `yaml:"Resources"`
}

// SAMResource ...
type SAMResource interface{}

// SAMFunction ...
type SAMFunction struct {
	SAMResource
	Type       string    `yaml:"Type"`
	Properties SAMFnProp `yaml:"Properties"`
}

// SAMFnProp ...
type SAMFnProp struct {
	Runtime string              `yaml:"Runtime"`
	Handler string              `yaml:"Handler"`
	CodeURI string              `yaml:"CodeUri"`
	Events  map[string]SAMEvent `yaml:"Events"`
}

// SAMEvent ...
type SAMEvent struct {
	Type       string  `yaml:"Type"`
	Properties SAMProp `yaml:"Properties"`
}

// SAMProp ...
type SAMProp struct {
	Path   string `yaml:"Path"`
	Method string `yaml:"Method"`
}
