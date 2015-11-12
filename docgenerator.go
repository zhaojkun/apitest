package testilla

import (
	"fmt"

	"github.com/alecthomas/jsonschema"
	"github.com/ghodss/yaml"
	"github.com/go-swagger/go-swagger/spec"
)

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
		s: seed,
	}
	gen.swagger.Swagger = "2.0" // from swagger doc: 'The value MUST be "2.0"'
	gen.swagger.Info = &spec.Info{}
	gen.swagger.Paths = &spec.Paths{Paths: map[string]spec.PathItem{}}

	return gen
}

// Generate implements IDocGenerator
func (g *swaggerYmlGenerator) Generate(tests []IApiTest) ([]byte, error) {
	for _, test := range tests {
		item := g.swagger.Paths.Paths[test.Path()]
		ops := map[string]spec.Operation{}
		for _, testCase := range test.TestCases() {
			// TODO: pregenerate and precollect operations and their parameters by some key first
			operation, err := generateOperation(testCase)
			if err != nil {
				return nil, err
			}

			ops[test.Method()] = operation
		}

		for method, op := range ops {
			switch method {
			case "GET":
				item.Get = &op
			case "POST":
				item.Post = &op
			case "PATCH":
				item.Patch = &op
			case "DELETE":
				item.Delete = &op
			case "PUT":
				item.Put = &op
			case "HEAD":
				item.Head = &op
			case "OPTIONS":
				item.Options = &op
			}
		}

		g.swagger.Paths.Paths[test.Path()] = item
	}

	return yaml.Marshal(g.s)
}

func generateOperation(testCase ApiTestCase) (op spec.Operation, err error) {
	op.Description = testCase.Description

	// generate parameters from Query Params
	for key, param := range testCase.QueryParams {
		specParam, err := generateSpecParam(key, param, "query")
		if err != nil {
			return op, err
		}
		op.Parameters = append(op.Parameters, specParam)
	}

	// generate parameters from Headers
	for key, param := range testCase.Headers {
		specParam, err := generateSpecParam(key, param, "header")
		if err != nil {
			return op, err
		}
		op.Parameters = append(op.Parameters, specParam)
	}

	// TODO: implement PATH parameters
	// for param := range testCase.PathParams {
	// 	specParam, err := generateSpecParam(param, "path")
	// 	if err != nil {
	// 		return op, err
	// 	}
	// 	op.Parameters = append(op.Parameters, specParam)
	// }
	if testCase.RequestBody != nil {
		specParam := spec.Parameter{}
		specParam.Name = "body"
		specParam.In = "body"
		specParam.Required = true
		specParam.Schema = generateSpecSchema(testCase.RequestBody)

		op.Parameters = append(op.Parameters, specParam)
	}
	return op, nil
}

func generateSpecParam(paramKey string, param ApiTestCaseParam, location string) (spec.Parameter, error) {
	specParam := spec.Parameter{}
	specParam.Name = paramKey
	specParam.In = location
	specParam.Required = param.Required
	specParam.Description = fmt.Sprintf("%v", param.Value)

	paramType, err := generateSpecSimpleType(param.Value)
	if err != nil {
		return specParam, fmt.Errorf("could not guess type of parameter '%s': %s", paramKey, err.Error())
	}
	specParam.Type = paramType

	return specParam, nil
}

func generateSpecSchema(item interface{}) *spec.Schema {
	schema := jsonschema.Reflect(item)
	return specSchemaFromJsonSchema(schema.Type)
}

func specSchemaFromJsonSchema(schema *jsonschema.Type) *spec.Schema {
	s := &spec.Schema{}
	s.Type = []string{schema.Type}
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
		s.Items.Schema = specSchemaFromJsonSchema(schema.Items)
	}

	if schema.Properties != nil {
		s.Properties = make(map[string]spec.Schema)
		for key, prop := range schema.Properties {
			s.Properties[key] = *specSchemaFromJsonSchema(prop)
		}
	}

	if schema.PatternProperties != nil {
		s.PatternProperties = make(map[string]spec.Schema)
		for key, prop := range schema.PatternProperties {
			s.PatternProperties[key] = *specSchemaFromJsonSchema(prop)
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
