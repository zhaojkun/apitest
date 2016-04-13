package apitest

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"

	// "github.com/advance512/yaml"
	"github.com/elgris/jsonschema"
	"github.com/go-raml/raml"
)

type ramlGenerator struct {
	doc raml.APIDefinition
}

func NewRamlGenerator(seed raml.APIDefinition) IDocGenerator {
	generator := &ramlGenerator{
		doc: seed,
	}

	generator.doc.RAMLVersion = "0.8"
	generator.doc.Resources = map[string]raml.Resource{}

	return generator
}

func (g *ramlGenerator) Generate(tests []IApiTest) ([]byte, error) {
	for _, test := range tests {
		// path MUST begin with '/'
		path := test.Path()
		if path[0] != '/' {
			path = "/" + path
		}

		resource := g.doc.Resources[path]
		if resource.UriParameters == nil {
			resource.UriParameters = map[string]raml.NamedParameter{}
		}

		m := raml.Method{}
		m.Responses = map[raml.HTTPCode]raml.Response{}
		m.Headers = map[raml.HTTPHeader]raml.Header{}
		m.QueryParameters = map[string]raml.NamedParameter{}
		m.Description = test.Description()
		m.Name = test.Method()

		processedHeaderParams := map[string]interface{}{}
		processedPathParams := map[string]interface{}{}
		processedQueryParams := map[string]interface{}{}

		for _, testCase := range test.TestCases() {
			m.Description = testCase.Description
			for key, param := range testCase.PathParams {
				if _, ok := processedPathParams[key]; ok {
					continue
				}

				uriParam := raml.NamedParameter{}
				uriParam.Name = key
				uriParam.Description = param.Description
				uriParam.Required = param.Required
				uriParam.Default = param.Value

				processedPathParams[key] = nil
				resource.UriParameters[key] = uriParam
			}

			for key, param := range testCase.Headers {
				if _, ok := processedHeaderParams[key]; ok {
					continue
				}

				h := raml.Header{}
				h.Name = key
				h.Description = param.Description
				h.Required = param.Required
				h.Default = param.Value

				m.Headers[raml.HTTPHeader(key)] = h
				processedHeaderParams[key] = nil
			}

			for key, param := range testCase.QueryParams {
				if _, ok := processedQueryParams[key]; ok {
					continue
				}

				uriParam := raml.NamedParameter{}
				uriParam.Name = key
				uriParam.Description = param.Description
				uriParam.Required = param.Required
				uriParam.Default = param.Value

				processedQueryParams[key] = nil
				m.QueryParameters[key] = uriParam
			}

			response := raml.Response{}
			response.Description = testCase.Description
			response.HTTPCode = raml.HTTPCode(testCase.ExpectedHttpCode)
			if testCase.ExpectedData != nil {
				schema := jsonschema.Reflect(testCase.ExpectedData)
				schemaBytes, _ := json.Marshal(schema)
				response.Bodies.DefaultSchema = string(schemaBytes)

				exampleBytes, _ := json.Marshal(testCase.ExpectedData)
				response.Bodies.DefaultExample = string(exampleBytes)
			}

			m.Responses[raml.HTTPCode(testCase.ExpectedHttpCode)] = response

		}

		// TODO: check if path has already assigned an method to some other test
		// return error if so
		switch test.Method() {
		case "GET":
			resource.Get = &m
		case "POST":
			resource.Post = &m
		case "PATCH":
			resource.Patch = &m
		case "DELETE":
			resource.Delete = &m
		case "PUT":
			resource.Put = &m
		case "HEAD":
			resource.Head = &m
		}

		g.doc.Resources[path] = resource
	}

	fmt.Printf("%+v\n\n\n", g.doc)
	d, e := yaml.Marshal(g.doc)
	fmt.Println("FOOOO", string(d))
	return d, e
}
