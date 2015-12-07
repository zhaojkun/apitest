package testilla

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/verdverm/frisby"
)

// ITestRunner is responsible for
// - receive a suite or suites
// - prepare some conditions to run suites (like setup frisby or whatsoever)
// - run tests of the suite
type ITestRunner interface {
	Run(tests []IApiTest, t *testing.T)
}

type basicRunner struct {
	DefaultHeaders map[string]string
	BaseUrl        string
}

func NewRunner(baseUrl string) *basicRunner {
	return &basicRunner{
		DefaultHeaders: make(map[string]string),
		BaseUrl:        baseUrl,
	}
}

func (r *basicRunner) Run(tests []IApiTest, t *testing.T) {
	for testIndex, test := range tests {
		for caseIndex, testCase := range test.TestCases() {

			err := r.runTest(testCase, test.Method(), test.Path())
			if err != nil {
				t.Errorf("error running test %d, case %d: %s", testIndex, caseIndex, err.Error())
			}
		}
	}
}

func (r *basicRunner) runTest(testCase ApiTestCase, method, path string) error {
	urlstring := r.BaseUrl + path
	url, err := testCase.Url(urlstring)
	if err != nil {
		return fmt.Errorf("could not prepare an url, '%s'", err.Error())
	}

	f := frisby.Create("")
	f.Method = method
	f.Url = url
	if testCase.RequestBody != nil {
		f.SetJson(testCase.RequestBody)
	}
	for name, param := range testCase.Headers {
		if stringValue, ok := param.Value.(string); ok {
			f.SetHeader(name, stringValue)
		} else {
			f.SetHeader(name, fmt.Sprintf("%v", param.Value))
		}
	}

	f.Send()
	if len(f.Errors()) > 0 {
		return f.Error()
	}

	// TODO: replace frisby with some basic HTTP client
	// asserting response code
	f.ExpectStatus(testCase.ExpectedHttpCode)

	// asserting headers
	if testCase.ExpectedHeaders != nil {
		for header, value := range testCase.ExpectedHeaders {

			f.ExpectHeader(header, value)
		}
	}

	if len(f.Errors()) > 0 {
		return f.Error()
	}

	// asserting body
	content, err := f.Resp.Content()
	if err != nil {
		return err
	}

	if testCase.ExpectedData != nil {
		// TODO 1: if expected string, do not convert it to JSON
		expectedData, err := objToJsonMap(testCase.ExpectedData)
		if err != nil {
			return fmt.Errorf("could not convert expected data to json map, '%s'", err.Error())
		}

		responseBody := map[string]interface{}{}
		if err := json.Unmarshal(content, &responseBody); err != nil {
			return err
		}

		if !reflect.DeepEqual(expectedData, responseBody) {
			expectedJson, _ := json.MarshalIndent(expectedData, "", "    ")
			actualJson, _ := json.MarshalIndent(responseBody, "", "    ")
			return fmt.Errorf("request and response are not equal.\nExpected: %s\nGot: %s", expectedJson, actualJson)
		}
	} else {
		if len(content) > 0 {
			return fmt.Errorf("expected empty response, got: %s", string(content))

		}
	}

	return nil
}

func objToJsonMap(obj interface{}) (map[string]interface{}, error) {
	js, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{}

	err = json.Unmarshal(js, &result)
	return result, err
}
