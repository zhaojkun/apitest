package main

import "github.com/seesawlabs/apitest"

type HelloTest struct {
}

func (t *HelloTest) Method() string      { return "GET" }
func (t *HelloTest) Description() string { return "Test for HelloWorld API handler" }
func (t *HelloTest) Path() string        { return "hello" }
func (t *HelloTest) TestCases() []apitest.ApiTestCase {
	return []apitest.ApiTestCase{
		{
			ExpectedHttpCode: 200,
			ExpectedData:     []byte("Hello World!\n"),
		},
	}
}
