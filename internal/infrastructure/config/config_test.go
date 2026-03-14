package config

import "testing"

func TestNewApplicationConfiguration(testingSuite *testing.T) {
	testCases := []struct {
		testName         string
		options          []Option
		expectationError bool
		portExpectation  int
	}{
		{testName: "default configuration", options: nil, expectationError: false, portExpectation: 8080},
		{testName: "override port", options: []Option{WithHTTPPort(9090)}, expectationError: false, portExpectation: 9090},
		{testName: "invalid port", options: []Option{WithHTTPPort(0)}, expectationError: true, portExpectation: 0},
	}

	for _, testCase := range testCases {
		testingSuite.Run(testCase.testName, func(testingContext *testing.T) {
			applicationConfiguration, configurationError := NewApplicationConfiguration(testCase.options...)
			if testCase.expectationError {
				if configurationError == nil {
					testingContext.Fatalf("expected configuration error but got nil")
				}
				return
			}
			if configurationError != nil {
				testingContext.Fatalf("unexpected configuration error: %v", configurationError)
			}
			if applicationConfiguration.HTTPPort != testCase.portExpectation {
				testingContext.Fatalf("unexpected port: got %d want %d", applicationConfiguration.HTTPPort, testCase.portExpectation)
			}
		})
	}
}
