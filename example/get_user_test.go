package main

import "github.com/seesawlabs/apitest"

type GetUserTest struct{}

func (t *GetUserTest) Method() string      { return "GET" }
func (t *GetUserTest) Description() string { return "Test for creating new user API" }
func (t *GetUserTest) Path() string        { return "/users/{user_id}" }
func (t *GetUserTest) TestCases() []apitest.ApiTestCase {
	return []apitest.ApiTestCase{
		{
			Description:      "User returned successfully",
			ExpectedHttpCode: 200,
			PathParams: apitest.ParamMap{
				"user_id": apitest.Param{Value: 2},
			},

			ExpectedData: User{
				ID:   2,
				Name: "Second User",
			},
		},
	}
}
