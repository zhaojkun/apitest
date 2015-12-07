package main

import (
	"github.com/octokit/go-octokit/octokit"

	"github.com/seesawlabs/testilla"
)

type GetUserTest struct {
}

func (t *GetUserTest) Method() string      { return "GET" }
func (t *GetUserTest) Description() string { return "Test for GetUser API handler" }
func (t *GetUserTest) Path() string        { return "user/{name}" }

func (t *GetUserTest) TestCases() []testilla.ApiTestCase {
	return []testilla.ApiTestCase{
		{
			Description: "Successful getting of user details",
			PathParams: testilla.ParamMap{
				"user": testilla.Param{Value: "elgris"},
			},

			ExpectedHttpCode: 200,
			ExpectedData:     octokit.User{},
		},
		{
			Description: "404 error in case user not found",
			PathParams: testilla.ParamMap{
				"user": testilla.Param{Value: "someveryunknown"},
			},

			ExpectedHttpCode: 404,
			ExpectedData:     "user 'someveryunknown' was not found",
		},
		{
			Description: "500 error in case something bad happens",
			PathParams: testilla.ParamMap{
				"user": testilla.Param{Value: "BadGuy"},
			},

			ExpectedHttpCode: 500,
			ExpectedData:     "BadGuy failed me :(",
		},
	}
}
