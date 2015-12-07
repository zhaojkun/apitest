package main

import (
	"testing"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/seesawlabs/testilla"
)

func TestApi(t *testing.T) {
	tests := []testilla.IApiTest{
		&HelloTest{},
		&GetUserTest{},
	}

	runner := testilla.NewRunner("http://127.0.0.1:1323/")
	runner.Run(tests, t)

	if !t.Failed() {
		seed := spec.Swagger{}
		seed.Host = "127.0.0.1:1323"
		seed.Produces = []string{"application/json"}
		seed.Consumes = []string{"application/json"}
		seed.Schemes = []string{"http"}
		seed.Info = &spec.Info{}
		seed.Info.Description = "Our very little example API with 2 endpoints"
		seed.Info.Title = "Example API"
		seed.Info.Version = "0.1"
		seed.BasePath = "/"

		generator := testilla.NewSwaggerYmlGenerator(seed)

		doc, err := generator.Generate(tests)
		if err != nil {
			t.Fatalf("could not generate docs: %s", err.Error())
		}

		t.Log(string(doc))
	}
}
