package apitest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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

func (r *basicRunner) Run(t *testing.T, tests ...IApiTest) {
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
			t.Logf("running test '%s'(%s), case %d", testName, test.Description(), caseIndex)
			r.runTest(t, testCase, test.Method(), test.Path())
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

func (r *basicRunner) encode(obj interface{}) ([]byte, error) {
	// TODO: make it configurable
	return json.Marshal(obj)
}

func (r *basicRunner) runTest(t *testing.T, testCase ApiTestCase, method, path string) {
	urlstring := r.BaseUrl + path
	url, err := testCase.Url(urlstring)
	if !assert.NoError(t, err, "could not prepare an url") {
		return
	}

	// TODO: prepare body
	var req *http.Request
	if testCase.RequestBody != nil {
		encoded, err := r.encode(testCase.RequestBody)
		if !assert.NoError(t, err, "could not encode body") {
			return
		}

		requestBody := bytes.NewBuffer(encoded)
		req, err = http.NewRequest(method, url, requestBody)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if !assert.NoError(t, err, "could not create HTTP request") {
		return
	}

	for name, value := range r.DefaultHeaders {
		req.Header.Set(name, value)
	}
	for name, param := range testCase.Headers {
		if stringValue, ok := param.Value.(string); ok {
			req.Header.Set(name, stringValue)
		} else {
			req.Header.Set(name, fmt.Sprintf("%v", param.Value))
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if !assert.NoError(t, err, "failed sending a request") {
		return
	}

	if !assert.NotNil(t, resp, "request to '%s' returned nil response", urlstring) {
		return
	}

	if !assert.Equal(t, testCase.ExpectedHttpCode, resp.StatusCode) {
		return
	}

	// asserting headers
	if testCase.ExpectedHeaders != nil {
		for header, value := range testCase.ExpectedHeaders {
			if !assert.Equal(t, value, resp.Header.Get(header)) {
				return
			}
		}
	}

	var responseBody []byte
	if resp.Body != nil {
		defer resp.Body.Close()

		responseBody, err = ioutil.ReadAll(resp.Body)
		if !assert.NoError(t, err) {
			return
		}
	}

	if testCase.ExpectedData != nil {
		expectedData := decodeExpected(testCase.ExpectedData)
		actualData := decodeResponse(responseBody)

		if !assert.Equal(t, actualData, expectedData, "request and response are not equal") {
			return
		}
	} else if !assert.Empty(t, responseBody, "expected empty response") {
		return
	}
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

	return string(data)
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
