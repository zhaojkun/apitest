package main

import "github.com/seesawlabs/apitest"

type CreateUserTest struct{}

func (t *CreateUserTest) Method() string      { return "POST" }
func (t *CreateUserTest) Description() string { return "Test for creating new user API" }
func (t *CreateUserTest) Path() string        { return "/users" }
func (t *CreateUserTest) TestCases() []apitest.ApiTestCase {
	return []apitest.ApiTestCase{
		{
			Description:      "User created successfully",
			ExpectedHttpCode: 201,
			Headers: apitest.ParamMap{
				"Content-Type": apitest.Param{Value: "application/json;charset=UTF-8"},
			},

			RequestBody: User{
				Name: "New User",
			},

			ExpectedData: User{
				ID:   3,
				Name: "New User",
			},
		},
	}
}
