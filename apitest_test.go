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

	// checking validity of generated swagger doc
	yamlMap := map[interface{}]interface{}{}
	err = yaml.Unmarshal(doc, &yamlMap)
	assert.NoError(t, err, "could not unmarshal generated doc into map")

	rawJSON, err := swag.YAMLToJSON(yamlMap)
	assert.NoError(t, err)

	swaggerDoc, err := spec.New(rawJSON, "")
	assert.NoError(t, err)

	err = validate.Spec(swaggerDoc, strfmt.Default)
	assert.NoError(t, err)

	// checking equality of generated and expected doc
	actual := map[interface{}]interface{}{}
	err = yaml.Unmarshal(doc, &actual)
	assert.NoError(t, err, "could not unmarshal generated doc into map")

	fixture, err := ioutil.ReadFile("fixtures/swagger/swagger.yml")
	assert.NoError(t, err, "could not read fixture file")

	expected := map[interface{}]interface{}{}
	err = yaml.Unmarshal(fixture, &expected)
	assert.NoError(t, err, "could not unmarshal fixture into map")

	assert.Equal(t, expected, actual)
}

func TestGenerateRaml(t *testing.T) {
	t.Skip("for now")
}

func getTests() []IApiTest {
	return []IApiTest{
		&HelloTest{},
		&GetUserTest{},
		&CreateUserTest{},
		&UpdateUserTest{},
		&DeleteUserTest{},
	}
}

func setupMock() {

	testUser := User{
		EventsURL:         "https://api.github.com/users/octocat/events{/privacy}",
		Followers:         20,
		FollowersURL:      "https://api.github.com/users/octocat/followers",
		Following:         0,
		FollowingURL:      "https://api.github.com/users/octocat/following{/other_user}",
		GistsURL:          "https://api.github.com/users/octocat/gists{/gist_id}",
		Hireable:          false,
		HTMLURL:           "https://github.com/octocat",
		Location:          "San Francisco",
		Login:             "octocat",
		Name:              "monalisa octocat",
		OrganizationsURL:  "https://api.github.com/users/octocat/orgs",
		PublicRepos:       2,
		ReceivedEventsURL: "https://api.github.com/users/octocat/received_events",
		ReposURL:          "https://api.github.com/users/octocat/repos",
		StarredURL:        "https://api.github.com/users/octocat/starred{/owner}{/repo}",
		SubscriptionsURL:  "https://api.github.com/users/octocat/subscriptions",
		Type:              "User",
		URL:               "https://api.github.com/users/octocat",
	}

	httpmock.RegisterResponder("GET", "http://testapi.my/hello",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, "Hello World!")
			return resp, nil
		},
	)

	httpmock.RegisterResponder("GET", "http://testapi.my/user/octocat",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, testUser)
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

	httpmock.RegisterResponder("POST", "http://testapi.my/user",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(201, testUser)
		},
	)

	httpmock.RegisterResponder("PATCH", "http://testapi.my/user/octocat",
		func(req *http.Request) (*http.Response, error) {
			patchedUser := testUser
			patchedUser.Name = "I Am Updated!"
			return httpmock.NewJsonResponse(200, patchedUser)
		},
	)

	httpmock.RegisterResponder("DELETE", "http://testapi.my/user/octocat",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewBytesResponse(204, nil)
			return resp, nil
		},
	)

	httpmock.RegisterResponder("DELETE", "http://testapi.my/user/someveryunknown",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(404, "user someveryunknown not found")
			return resp, nil
		},
	)

	httpmock.RegisterResponder("DELETE", "http://testapi.my/user/BadGuy",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(500, "BadGuy failed me :(")
			return resp, nil
		},
	)
}
