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
	for _, test := range tests {
		testName := extractTestName(test)
		// setup test
		if setuppable, ok := test.(ISetuppable); ok {
			t.Logf("setting up test '%s'(%s)...", testName, test.Description())

			if err := setuppable.SetUp(); err != nil {
				t.Errorf("error setting up test '%s'(%s): %s",
					testName, test.Description(), err.Error())

				continue
			}
		}

		// run test
		for caseIndex, testCase := range test.TestCases() {

			err := r.runTest(testCase, test.Method(), test.Path())
			if err != nil {
				t.Errorf("error running test '%s'(%s), case %d: %s",
					testName, test.Description(), caseIndex, err.Error())
			}
		}

		// teardown test
		if teardownable, ok := test.(ITeardownable); ok {
			t.Logf("tearing down test '%s'(%s)...", testName, test.Description())

			if err := teardownable.TearDown(); err != nil {
				t.Errorf("error cleaning up after a test '%s'(%s): %s",
					testName, test.Description(), err.Error())
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
	f.SetHeaders(r.DefaultHeaders)

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
		expectedData := decodeExpected(testCase.ExpectedData)
		actualData := decodeResponse(content)

		if !reflect.DeepEqual(expectedData, actualData) {
			expectedJson, _ := json.MarshalIndent(expectedData, "", "    ")
			actualJson, _ := json.MarshalIndent(actualData, "", "    ")
			return fmt.Errorf("request and response are not equal.\nExpected: %s\nGot: %s", expectedJson, actualJson)
		}
	} else {
		if len(content) > 0 {
			return fmt.Errorf("expected empty response, got: %s", string(content))

		}
	}

	return nil
}

// decodeExpected processes expected data into representation used for comparison to actual data
// (converts object to map, does some type conversion)
func decodeExpected(data interface{}) interface{} {
	if decoded, err := objToJsonMap(data); err == nil {
		return decoded
	}

	return data
}

// decodeResponse processes response data into representation used for comparison to expected data
func decodeResponse(data []byte) interface{} {
	var dataMap map[string]interface{}
	if err := json.Unmarshal(data, &dataMap); err == nil {
		return dataMap
	}

	return data
}

func extractTestName(value interface{}) string {
	return reflect.TypeOf(value).String()
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
