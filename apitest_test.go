package apitest

import (
	"io/ioutil"
	"net/http"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/go-swagger/go-swagger/strfmt"
	"github.com/go-swagger/go-swagger/swag"
	"github.com/go-swagger/go-swagger/validate"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestRunApi(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	setupMock()

	tests := getTests()

	runner := NewRunner("http://testapi.my")
	runner.Run(t, tests...)
}

func TestGenerateSwaggerYAML(t *testing.T) {
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

	generator := NewSwaggerGeneratorYAML(seed)
	tests := getTests()

	doc, err := generator.Generate(tests)
	assert.NoError(t, err, "could not generate docs")

	actual := map[string]interface{}{}
	err = yaml.Unmarshal(doc, &actual)
	assert.NoError(t, err, "could not unmarshal generated doc into map")

	fixture, err := ioutil.ReadFile("fixtures/swagger/swagger.yml")
	assert.NoError(t, err, "could not read fixture file")

	expected := map[string]interface{}{}
	err = yaml.Unmarshal(fixture, &expected)
	assert.NoError(t, err, "could not unmarshal fixture into map")

	assert.Equal(t, expected, actual)
}

func TestValidateSwaggerYAML(t *testing.T) {
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

	generator := NewSwaggerGeneratorYAML(seed)

	tests := getTests()

	doc, err := generator.Generate(tests)
	assert.NoError(t, err, "could not generate docs")

	yamlMap := map[interface{}]interface{}{}
	err = yaml.Unmarshal(doc, &yamlMap)
	assert.NoError(t, err, "could not unmarshal generated doc into map")

	rawJSON, err := swag.YAMLToJSON(yamlMap)
	assert.NoError(t, err)

	swaggerDoc, err := spec.New(rawJSON, "")
	assert.NoError(t, err)

	err = validate.Spec(swaggerDoc, strfmt.Default)
	assert.NoError(t, err)
}

func getTests() []IApiTest {
	return []IApiTest{
		&HelloTest{},
		&GetUserTest{},
	}
}

func TestGenerateRaml(t *testing.T) {
	t.Skip("for now")
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
			resp := httpmock.NewStringResponse(200, content)
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
