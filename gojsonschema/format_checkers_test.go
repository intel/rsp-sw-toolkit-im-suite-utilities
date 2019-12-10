/* Apache v2 license
*  Copyright (C) <2019> Intel Corporation
*
*  SPDX-License-Identifier: Apache-2.0
 */

package gojsonschema

import (
	"testing"
)

type testCase struct {
	name     string
	input    string
	expected bool
}

func TestEmailFormatCheckers(t *testing.T) {
	emailFormatChecker := EmailFormatChecker{}
	testCases := []testCase{
		{name: "normal ASCII characters for email", input: "abc@intel.com", expected: true},
		{name: "normal numerical characters for email", input: "123@intel.com", expected: true},
		{name: "normal ASCII mixed with numerical characters for email", input: "Abc123@intel.com", expected: true},
		{name: "some special characters for email", input: "#!+.^&%`@intel.com", expected: true},
		{name: "bad multiple @ for email", input: "bad@test@intel.com", expected: false},
		{name: "empty before @ for email", input: "@intel.com", expected: false},
		{name: "empty after @ for email", input: "test@", expected: false},
		{name: "special . after @ for email", input: "test@..", expected: false},
		{name: "invalid \" character for email", input: "\"test@inte.com", expected: false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if ok := emailFormatChecker.IsFormat(testCase.input); ok != testCase.expected {
				t.Errorf("Test name [%s] for email %s: expecting %v but found %v",
					testCase.name, testCase.input, testCase.expected, ok)
			}
		})
	}
}
