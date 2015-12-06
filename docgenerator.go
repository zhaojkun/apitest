package testilla

import (
	"fmt"

	"github.com/elgris/jsonschema"
	"github.com/ghodss/yaml"
	"github.com/go-swagger/go-swagger/spec"
)

// ITaggable is an interface that can tell doc generator
// that some test provides a tag. Usable for swagger documentation
// where tags help to group API endpoints
type ITaggable interface {
	Tag() string
}

// IDocGenerator describes a generator of documentation
// uses test suite as a source of information about API endpoints
type IDocGenerator interface {
	Generate(tests []IApiTest) ([]byte, error)
}

type swaggerYmlGenerator struct {
	swagger spec.Swagger
}

// NewSwaggerYmlGenerator initializes new generator with initial swagger spec
// as a seed. The generator produces YML output
func NewSwaggerYmlGenerator(seed spec.Swagger) IDocGenerator {
	gen := &swaggerYmlGenerator{
		swagger: seed,
	}
	gen.swagger.Swagger = "2.0" // from swagger doc: 'The value MUST be "2.0"'
	gen.swagger.Paths = &spec.Paths{Paths: map[string]spec.PathItem{}}

	return gen
}

// Generate implements IDocGenerator
func (g *swaggerYmlGenerator) Generate(tests []IApiTest) ([]byte, error) {
	g.swagger.Definitions = spec.Definitions{}

	for _, test := range tests {
		path := g.swagger.Paths.Paths[test.Path()] // TODO: 2 tests on the same API with the same response code conflict
		op, err := g.generateSwaggerOperation(test, g.swagger.Definitions)
		if err != nil {
			return nil, err
		}

		// TODO: check if path has already assigned an operation to some other test
		// return error if so
		switch test.Method() {
		case "GET":
			path.Get = &op
		case "POST":
			path.Post = &op
		case "PATCH":
			path.Patch = &op
		case "DELETE":
			path.Delete = &op
		case "PUT":
			path.Put = &op
		case "HEAD":
			path.Head = &op
		case "OPTIONS":
			path.Options = &op
		}

		g.swagger.Paths.Paths[test.Path()] = path
	}

	return yaml.Marshal(g.swagger)
}

func (g *swaggerYmlGenerator) generateSwaggerOperation(test IApiTest, defs spec.Definitions) (spec.Operation, error) {

	op := spec.Operation{}
	op.Responses = &spec.Responses{}
	op.Responses.StatusCodeResponses = map[int]spec.Response{}

	var description string
	params := map[string]spec.Parameter{}
	for _, testCase := range test.TestCases() {
		// parameter definitions are collected from 2xx tests only
		if testCase.ExpectedHttpCode >= 200 && testCase.ExpectedHttpCode < 300 {
			description = testCase.Description

			for key, param := range testCase.QueryParams {
				specParam, err := generateSpecParam(key, param, "query")
				if err != nil {
					return op, err
				}

				params[specParam.Name+specParam.In] = specParam
			}
			for key, param := range testCase.PathParams {
				specParam, err := generateSpecParam(key, param, "path")
				if err != nil {
					return op, err
				}

				params[specParam.Name+specParam.In] = specParam
			}

			for key, param := range testCase.Headers {
				specParam, err := generateSpecParam(key, param, "header")
				if err != nil {
					return op, err
				}

				params[specParam.Name+specParam.In] = specParam
			}

			if testCase.RequestBody != nil {
				specParam := spec.Parameter{}
				specParam.Name = "body"
				specParam.In = "body"
				specParam.Required = true
				specParam.Default = testCase.RequestBody

				specParam.Schema = generateSpecSchema(testCase.RequestBody, defs)

				params[specParam.Name+specParam.In] = specParam
			}
		}

		response := spec.Response{}
		response.Description = testCase.Description
		if testCase.ExpectedData != nil {
			response.Schema = generateSpecSchema(testCase.ExpectedData, defs)
			response.Examples = map[string]interface{}{
				"application/json": testCase.ExpectedData,
			}
		}

		op.Responses.StatusCodeResponses[testCase.ExpectedHttpCode] = response
	}

	for _, param := range params {
		op.Parameters = append(op.Parameters, param)
	}
	op.Summary = description
	if taggable, ok := test.(ITaggable); ok {
		op.Tags = []string{taggable.Tag()}
	}

	return op, nil
}

func generateSpecParam(paramKey string, param Param, location string) (spec.Parameter, error) {
	specParam := spec.Parameter{}
	specParam.Name = paramKey
	specParam.In = location
	specParam.Required = param.Required
	specParam.Description = param.Description
	specParam.Default = param.Value

	paramType, err := generateSpecSimpleType(param.Value)
	if err != nil {
		return specParam, fmt.Errorf("could not guess type of parameter '%s': %s", paramKey, err.Error())
	}
	specParam.Type = paramType

	return specParam, nil
}

func generateSpecSchema(item interface{}, defs spec.Definitions) *spec.Schema {
	refl := jsonschema.Reflect(item)
	schema := specSchemaFromJsonType(refl.Type)

	schema.Definitions = map[string]spec.Schema{}
	for name, def := range refl.Definitions {
		defs[name] = *specSchemaFromJsonType(def)
	}

	return schema
}

func specSchemaFromJsonType(schema *jsonschema.Type) *spec.Schema {
	s := &spec.Schema{}
	if schema.Type != "" {
		s.Type = []string{schema.Type}
	}
	if schema.Ref != "" {
		s.Ref = spec.MustCreateRef(schema.Ref)
	}

	s.Format = schema.Format
	s.Required = schema.Required

	// currently there is no way to determine whether there is MaxLength or MinLength
	// defined. Need to fix jsonschema library and switch type from int to *int
	// s.MaxLength = schema.MaxLength
	// s.MinLength = schema.MinLength
	s.Pattern = schema.Pattern
	s.Enum = schema.Enum
	s.Default = schema.Default
	s.Title = schema.Title
	s.Description = schema.Description

	if schema.Items != nil {
		s.Items = &spec.SchemaOrArray{}
		s.Items.Schema = specSchemaFromJsonType(schema.Items)
	}

	if schema.Properties != nil {
		s.Properties = make(map[string]spec.Schema)
		for key, prop := range schema.Properties {
			s.Properties[key] = *specSchemaFromJsonType(prop)
		}
	}

	if schema.PatternProperties != nil {
		s.PatternProperties = make(map[string]spec.Schema)
		for key, prop := range schema.PatternProperties {
			s.PatternProperties[key] = *specSchemaFromJsonType(prop)
		}
	}

	switch string(schema.AdditionalProperties) {
	case "true":
		s.AdditionalProperties = &spec.SchemaOrBool{Allows: true}
	case "false":
		s.AdditionalProperties = &spec.SchemaOrBool{Allows: false}
	}

	return s
}

func generateSpecSimpleType(value interface{}) (string, error) {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return "integer", nil
	case float32, float64:
		return "number", nil
	case bool:
		return "boolean", nil
	case string:
		return "string", nil
	}

	return "", fmt.Errorf("value of complex type '%T' provided, simple type expected", value)
}
