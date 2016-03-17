package main

import "github.com/seesawlabs/apitest"

type DeleteUserTest struct{}

func (t *DeleteUserTest) Method() string      { return "DELETE" }
func (t *DeleteUserTest) Description() string { return "Test for creating new user API" }
func (t *DeleteUserTest) Path() string        { return "/users/{user_id}" }
func (t *DeleteUserTest) TestCases() []apitest.ApiTestCase {
	return []apitest.ApiTestCase{
		{
			Description:      "User deleted successfully",
			ExpectedHttpCode: 204,
			PathParams: apitest.ParamMap{
				"user_id": apitest.Param{Value: 3},
			},
		},
	}
}
