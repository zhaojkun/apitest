package apitest

type ramlGenerator struct {
}

func NewRamlGenerator() IDocGenerator {
	gen := &ramlGenerator{}

	return gen
}

type APIDefinition struct {
	// TODO: hardcode
	RAMLVersion string `yaml:"raml_version"`
	// TODO: from seed
	Title string `yaml:"title"`
	// TODO: from seed
	Version string `yaml:"version"`
	// TODO: from seed
	BaseUri string `yaml:"baseUri"`
	// TODO: from seed
	BaseUriParameters map[string]NamedParameter `yaml:"baseUriParameters"`

	UriParameters map[string]NamedParameter `yaml:"uriParameters"`
	// TODO: from seed
	Protocols []string `yaml:"protocols"`
	// TODO: from seed
	Description     string                      `yaml:"description"`
	MediaType       string                      `yaml:"mediaType"`
	SecuritySchemes []map[string]SecurityScheme `yaml:"securitySchemes"`
	Types           []map[string]string
	ResourceTypes   []map[string]ResourceType `yaml:"resourceTypes"`
	Traits          []map[string]Trait        `yaml:"traits"`
	SecuredBy       []DefinitionChoice        `yaml:"securedBy"`
	Documentation   []Documentation           `yaml:"documentation"`
	// Annotations
	// AnnotationTypes
	// Uses
	Resources map[string]Resource `yaml:",regexp:/.*"`
}

func (g *ramlGenerator) Generate(tests []IApiTest) ([]byte, error) {

}
