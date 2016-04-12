package apitest

import "github.com/seesawlabs/raml"

type ramlGenerator struct {
	doc raml.APIDefinition
}

func NewRamlGenerator(seed raml.APIDefinition) IDocGenerator {
	generator := &ramlGenerator{
		doc: seed,
	}

	generator.doc.RAMLVersion = "1.0"

	return generator
}

func (g *ramlGenerator) Generate(tests []IApiTest) ([]byte, error) {

}
