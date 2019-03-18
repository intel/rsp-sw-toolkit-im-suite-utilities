/*
 * INTEL CONFIDENTIAL
 * Copyright (2018) Intel Corporation.
 *
 * The source code contained or described herein and all documents related to the source code ("Material")
 * are owned by Intel Corporation or its suppliers or licensors. Title to the Material remains with
 * Intel Corporation or its suppliers and licensors. The Material may contain trade secrets and proprietary
 * and confidential information of Intel Corporation and its suppliers and licensors, and is protected by
 * worldwide copyright and trade secret laws and treaty provisions. No part of the Material may be used,
 * copied, reproduced, modified, published, uploaded, posted, transmitted, distributed, or disclosed in
 * any way without Intel/'s prior express written permission.
 * No license under any patent, copyright, trade secret or other intellectual property right is granted
 * to or conferred upon you by disclosure or delivery of the Materials, either expressly, by implication,
 * inducement, estoppel or otherwise. Any license under such intellectual property rights must be express
 * and approved by Intel in writing.
 * Unless otherwise agreed by Intel in writing, you may not remove or alter this notice or any other
 * notice embedded in Materials by Intel or Intel's suppliers or licensors in any way.
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
