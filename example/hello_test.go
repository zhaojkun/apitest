package main

import "github.com/seesawlabs/testilla"

type HelloTest struct {
}

func (t *HelloTest) Method() string      { return "GET" }
func (t *HelloTest) Description() string { return "Test for HelloWorld API handler" }
func (t *HelloTest) Path() string        { return "hello" }
func (t *HelloTest) TestCases() []testilla.ApiTestCase {
	return []testilla.ApiTestCase{
		{
			ExpectedHttpCode: 200,
			ExpectedData:     "Hello world!",
		},
	}
}
