package main

import "github.com/seesawlabs/apitest"

type UpdateUserTest struct{}

func (t *UpdateUserTest) Method() string      { return "PATCH" }
func (t *UpdateUserTest) Description() string { return "Test for creating new user API" }
func (t *UpdateUserTest) Path() string        { return "/users/{user_id}" }
func (t *UpdateUserTest) TestCases() []apitest.ApiTestCase {
	return []apitest.ApiTestCase{
		{
			Description:      "User updated successfully",
			ExpectedHttpCode: 200,
			Headers: apitest.ParamMap{
				"Content-Type": apitest.Param{Value: "application/json;charset=UTF-8"},
			},
			PathParams: apitest.ParamMap{
				"user_id": apitest.Param{Value: 1},
			},

			RequestBody: User{
				Name: "I Am Updated!",
			},

			ExpectedData: User{
				ID:   1,
				Name: "I Am Updated!",
			},
		},
	}
}
