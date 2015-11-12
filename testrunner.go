package testilla

import "testing"

// ITestRunner is responsible for
// - receive a suite or suites
// - prepare some conditions to run suites (like setup frisby or whatsoever)
// - run tests of the suite
type ITestRunner interface {
	Run(tests []IApiTest, t *testing.T)
}
