package main

import (
	"fmt"
	"testing"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/seesawlabs/apitest"
)

func TestRunApi(t *testing.T) {

	tests := []apitest.IApiTest{
		&GetUserTest{},
		&UpdateUserTest{},
		&CreateUserTest{},
		&DeleteUserTest{},
	}

	runner := apitest.NewRunner("http://localhost:1323")
	runner.Run(t, tests...)

	if !t.Failed() {
		generateSwaggerYAML(t, tests)
	}
}

func generateSwaggerYAML(t *testing.T, tests []apitest.IApiTest) {
	seed := spec.Swagger{}
	seed.Host = "localhost"
	seed.Produces = []string{"application/json"}
	seed.Consumes = []string{"application/json"}
	seed.Schemes = []string{"http"}
	seed.Info = &spec.Info{}
	seed.Info.Description = "Example API"
	seed.Info.Title = "Example API"
	seed.Info.Version = "0.1"
	seed.BasePath = "/"

	generator := apitest.NewSwaggerGeneratorYAML(seed)

	doc, err := generator.Generate(tests)
	if err != nil {
		t.Fatalf("could not generate doc: %s", err.Error())
	}

	fmt.Println(string(doc))
}
