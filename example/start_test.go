package main

import (
	"net/http"
	"testing"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/jarcoal/httpmock"
	"github.com/seesawlabs/apitest"
)

func TestApi(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	setupMock()

	tests := []apitest.IApiTest{
		&HelloTest{},
		&GetUserTest{},
	}

	runner := apitest.NewRunner("http://testapi.my/")
	runner.Run(t, tests...)

	if !t.Failed() {
		seed := spec.Swagger{}
		seed.Host = "testapi.my"
		seed.Produces = []string{"application/json"}
		seed.Consumes = []string{"application/json"}
		seed.Schemes = []string{"http"}
		seed.Info = &spec.Info{}
		seed.Info.Description = "Our very little example API with 2 endpoints"
		seed.Info.Title = "Example API"
		seed.Info.Version = "0.1"
		seed.BasePath = "/"

		generator := apitest.NewSwaggerYmlGenerator(seed)

		doc, err := generator.Generate(tests)
		if err != nil {
			t.Fatalf("could not generate docs: %s", err.Error())
		}

		t.Log(string(doc))
	}
}

func setupMock() {

	httpmock.RegisterResponder("GET", "http://testapi.my/hello",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, "Hello World!")
			return resp, nil
		},
	)

	httpmock.RegisterResponder("GET", "http://testapi.my/user/octocat",
		func(req *http.Request) (*http.Response, error) {
			content := `{
  "login": "octocat",
  "url": "https://api.github.com/users/octocat",
  "html_url": "https://github.com/octocat",
  "followers_url": "https://api.github.com/users/octocat/followers",
  "following_url": "https://api.github.com/users/octocat/following{/other_user}",
  "gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
  "starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
  "subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
  "organizations_url": "https://api.github.com/users/octocat/orgs",
  "repos_url": "https://api.github.com/users/octocat/repos",
  "events_url": "https://api.github.com/users/octocat/events{/privacy}",
  "received_events_url": "https://api.github.com/users/octocat/received_events",
  "type": "User",
  "name": "monalisa octocat",
  "location": "San Francisco",
  "public_repos": 2,
  "followers": 20
}`
			resp, _ := httpmock.NewJsonResponse(200, content)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("GET", "http://testapi.my/user/someveryunknown",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(404, "user someveryunknown not found")
			return resp, nil
		},
	)

	httpmock.RegisterResponder("GET", "http://testapi.my/user/BadGuy",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(500, "BadGuy failed me :(")
			return resp, nil
		},
	)
}
